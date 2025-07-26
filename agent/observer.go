package agent

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskProcessing
	TaskCompleted
	TaskFailed
)

func (ts TaskStatus) String() string {
	switch ts {
	case TaskPending:
		return "pending"
	case TaskProcessing:
		return "processing"
	case TaskCompleted:
		return "completed"
	case TaskFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type TaskObserver struct {
	TaskID    string
	Status    TaskStatus
	Progress  int
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ObserverManager struct {
	observers map[string]*TaskObserver
	mutex     sync.RWMutex
}

var observerManager = &ObserverManager{
	observers: make(map[string]*TaskObserver),
}

func GetObserverManager() *ObserverManager {
	return observerManager
}

func (om *ObserverManager) RegisterTask(taskID string) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	om.observers[taskID] = &TaskObserver{
		TaskID:    taskID,
		Status:    TaskPending,
		Progress:  0,
		Message:   "Task created",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Printf("Task registered: %s", taskID)
}

func (om *ObserverManager) UpdateTask(taskID string, status TaskStatus, progress int, message string) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	if observer, exists := om.observers[taskID]; exists {
		observer.Status = status
		observer.Progress = progress
		observer.Message = message
		observer.UpdatedAt = time.Now()

		log.Printf("Task updated: %s - %s (%d%%): %s", taskID, status, progress, message)
	}
}

func (om *ObserverManager) GetTask(taskID string) (*TaskObserver, bool) {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	observer, exists := om.observers[taskID]
	return observer, exists
}

func (om *ObserverManager) RemoveTask(taskID string) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	delete(om.observers, taskID)
	log.Printf("Task removed: %s", taskID)
}

func (om *ObserverManager) ListTasks() map[string]*TaskObserver {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	// Return a copy to avoid race conditions
	tasks := make(map[string]*TaskObserver)
	for k, v := range om.observers {
		tasks[k] = v
	}

	return tasks
}

// Helper function to update task progress during video processing
func UpdateTaskProgress(taskID string, step string, progress int) {
	message := fmt.Sprintf("Processing: %s", step)
	observerManager.UpdateTask(taskID, TaskProcessing, progress, message)
}
