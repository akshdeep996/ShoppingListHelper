package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	AuthorId  string     `json:"authorId"`
	TaskItems []TaskItem `json:"taskItems"`
	Progress  float64    `json:"progress"`
}

type TaskItem struct {
	ID         string `json:"id"`
	TaskDetail string `json:"taskDetail"`
	IsDone     bool   `json:"isDone"`
	Quantity   int    `json:"quantity"`
	AuthorId   string `json:"authorId"`
}

var taskArray = make([]Task, 0)

func print() {
	fmt.Println("=================")
	fmt.Println("Task List")
	fmt.Println("=================")
	for _, task := range taskArray {
		fmt.Println(task)
		i := 0
		for i < len(task.TaskItems) {
			fmt.Println(task.TaskItems[i])
			i++
		}
	}
	fmt.Println("=================")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func createTask(title string, author string, taskItems []TaskItem) Task {
	taskObject := Task{
		ID:        (uuid.New()).String(),
		Title:     title,
		AuthorId:  author,
		TaskItems: taskItems,
	}
	for id, task := range taskObject.TaskItems {
		task.ID = (uuid.New()).String()
		taskObject.TaskItems[id] = task
	}
	taskObject = updateProgress(taskObject)
	return taskObject
}

func updateProgress(taskObject Task) Task {
	completedTaskCount := 0
	for _, task := range taskObject.TaskItems {
		if task.IsDone == true {
			completedTaskCount++
		}
	}
	if len(taskObject.TaskItems) == 0 {
		taskObject.Progress = 0
	} else {
		taskObject.Progress = float64(completedTaskCount) / float64(len(taskObject.TaskItems)) * 100.0
	}
	return taskObject
}

func addTask(task Task) {
	taskArray = append(taskArray, task)
}

func addTaskItem(taskListID string, taskItem TaskItem) Task {
	for id, task := range taskArray {
		if task.ID == taskListID {
			task.TaskItems = append(task.TaskItems, taskItem)
			taskArray[id] = task
			return task
		}
	}
	return Task{}
}

func getTaskList() []Task {

	return taskArray
}

func getTaskListByID(taskID string) (Task, error) {
	for _, task := range taskArray {
		if task.ID == taskID {
			return task, nil
		}
	}
	return Task{}, errors.New("Task not found")
}

func deleteTaskListByID(taskID string) (bool, error) {
	i := 0
	isFound := false
	for _, task := range taskArray {
		if task.ID == taskID {
			isFound = true
			break
		}
		i++
	}

	if isFound {
		taskArray[i] = taskArray[len(taskArray)-1]
		taskArray = taskArray[:len(taskArray)-1]
		return true, nil
	} else {
		return false, errors.New("Task not found")
	}

}

func updateTask(newTask Task) Task {
	newTask = updateProgress(newTask)
	for i, task := range taskArray {
		if task.ID == newTask.ID {
			taskArray[i] = newTask
			break
		}
	}
	return newTask
}

// APIS
func getTaskListAPI(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, gin.H{"Access-Control-Allow-Origin": "*"})
	// c.IndentedJSON(http.StatusOK, gin.H{"Access-Control-Allow-Headers": "authentication, content-type"})
	// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	c.IndentedJSON(http.StatusOK, getTaskList())

}

func getTaskListByIDAPI(c *gin.Context) {
	taskID := c.Param("taskID")
	task, err := getTaskListByID(taskID)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func deleteTaskListByIDAPI(c *gin.Context) {
	taskID := c.Param("taskID")
	status, err := deleteTaskListByID(taskID)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	if !status {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Delete operation failed"})
		return
	}
	// c.IndentedJSON(http.StatusOK, gin.H{"Access-Control-Allow-Origin": "*"})
	// c.IndentedJSON(http.StatusOK, gin.H{"Access-Control-Allow-Headers": "authentication, content-type"})

	c.IndentedJSON(http.StatusOK, status)
}

func createTaskListAPI(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		return
	}
	newTask := createTask(task.Title, task.AuthorId, task.TaskItems)
	addTask(newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func putTaskListAPI(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		return
	}
	task = updateTask(task)
	c.IndentedJSON(http.StatusCreated, task)
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/tasks", getTaskListAPI)
	router.GET("/tasks/:taskID", getTaskListByIDAPI)
	router.PUT("/tasks", putTaskListAPI)
	router.POST("/tasks", createTaskListAPI)
	router.DELETE("/tasks/:taskID", deleteTaskListByIDAPI)

	router.Run("localhost:9090")

	// Uncomment to test he changes locally by printing messaged to console
	localTest()
}

func localTest() {
	task := createTask(
		"Shopping List",
		"Akshdeep Singh",
		make([]TaskItem, 0),
	)
	addTask(task)

	addTaskItem(task.ID, TaskItem{
		TaskDetail: "Potato 2 Kg",
		AuthorId:   "Akshdeep Rajawat",
	})
	addTaskItem(task.ID, TaskItem{
		TaskDetail: "Sugar 1 Kg",
		AuthorId:   "Geetha Ghulekar",
	})

	fmt.Println(getTaskList())
	fmt.Println(getTaskListByID(task.ID))
	print()
	task2 := createTask(
		"Shopping List - Updated",
		"Random User",
		make([]TaskItem, 0),
	)
	task2.ID = task.ID

	updateTask(task2)

	fmt.Println(getTaskList())
	fmt.Println(getTaskListByID(task.ID))
	print()
}
