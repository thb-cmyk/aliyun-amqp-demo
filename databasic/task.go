package databasic

import (
	"context"
	"errors"
	"time"
)

type TaskNode struct {
	Id string

	List *ListNode /* it is a continer that is used to orgnize the parent type as a list */

	Method *ProceNode /* it is a method to process the raw data */

	Raw_list *ListNode /* it is a list holding raw data, which receive from the global channel */
	Raw_max  int       /* the memeber is unused */
	Raw_num  int       /* the number of Raw_list hold raw data amount */

	Timepeice time.Duration /* Timepeice is the max value of running time each calling. */

	Timeout        time.Time /* Timeout is a timepeice. The Timeout is set, while the Raw_list is empty */
	Tiemout_is_set bool      /* the member is used to check wether the timeout is set or not. */

	Cancel    chan bool
	Goroutine bool
}

func TaskNode_register(id string, method *ProceNode, timepeice time.Duration) *TaskNode {
	if id == "" || method == nil || timepeice <= 0 {
		return nil
	}

	rawnode := RawNode_create("sentry", nil)

	tasknode := new(TaskNode)
	tasknode.Id = id
	tasknode.Method = method
	tasknode.Timepeice = time.Duration(timepeice)
	tasknode.List = ListNode_create(tasknode)
	tasknode.Raw_list = ListNode_create(rawnode)
	tasknode.Raw_max = 50
	tasknode.Raw_num = 0
	tasknode.Cancel = make(chan bool)
	tasknode.Goroutine = false

	/* add the dataclass to the global list */
	ok := ListNode_insert_next(global_tasknode_entry, tasknode.List)
	if !ok {
		tasknode.List.Parent = nil
		tasknode.Raw_list.Parent = nil
		tasknode.List = nil
		tasknode.Raw_list = nil
		return nil
	}
	global_tasknode_num++

	return tasknode

}

func TaskNode_find(id string) *TaskNode {

	listnode := global_tasknode_entry

	if listnode.Next == listnode || listnode.Prev == listnode {
		return nil
	}

	if id == "last" {
		return listnode.Prev.Parent.(*TaskNode)
	} else if id == "first" {
		return listnode.Next.Parent.(*TaskNode)
	}

	for i := 0; i < global_tasknode_num; i++ {
		listnode = ListNode_index_node(listnode, 1)
		parent := listnode.Parent.(*TaskNode)
		if parent.Id == id {
			return parent
		}
	}
	return nil
}

func (tn *TaskNode) TaskNode_unregister() bool {

	ok := ListNode_delete(tn.List)
	if !ok {
		return false
	}
	close(tn.Cancel)
	global_tasknode_num--

	return true
}

func (tn *TaskNode) TaskNode_add(rawnode *RawNode) bool {

	if rawnode == nil {
		return false
	}
	data_entry := tn.Raw_list
	listnode := rawnode.List
	ok := ListNode_insert_next(data_entry, listnode)
	if !ok {

		return false
	}
	tn.Raw_num++

	return true
}

func (tn *TaskNode) TaskNode_search(id string) *RawNode {

	listnode := tn.Raw_list

	if listnode.Next == listnode || listnode.Prev == listnode {
		return nil
	}

	if id == "last" {
		return listnode.Prev.Parent.(*RawNode)
	} else if id == "first" {
		return listnode.Next.Parent.(*RawNode)
	}

	for index := 0; index < tn.Raw_num; index++ {
		listnode = ListNode_index_node(listnode, 1)
		parent := listnode.Parent.(*RawNode)
		if parent.Id == id {
			return parent
		}
	}

	return nil
}

func (tn *TaskNode) TaskNode_remove(rawnode *RawNode) bool {

	if rawnode == nil {

		return false
	} else {
		ok := ListNode_delete(rawnode.List)
		if !ok {

			return false
		}
		tn.Raw_num--
	}

	return true
}

func (tn *TaskNode) TaskNode_update_mthod(method *ProceNode, ctx context.Context) bool {

	if method == nil {

		return false
	} else {
		tn.Method = method

		return true
	}
}

func (tn *TaskNode) TaskNode_update_timepeice(timepeice int64) error {

	if timepeice < 0 {
		return errors.New("argument error")
	} else {
		tn.Timepeice = time.Duration(timepeice)
		return nil
	}
}

func (tn *TaskNode) TaskNode_set_timeout(timeout time.Duration) error {

	if timeout < 0 {
		return errors.New("argument error")
	}

	duration := timeout
	tn.Timeout = <-time.After(duration)
	tn.Tiemout_is_set = true

	return nil
}

func (tn *TaskNode) TaskNode_unset_timeout() {

	tn.Tiemout_is_set = false
	tn.Timeout = time.Now()

}

func (tn *TaskNode) TaskNode_is_timeout() bool {
	if tn.Tiemout_is_set {
		if tn.Timeout.Before(time.Now()) {
			return false
		}
	}
	return true
}
