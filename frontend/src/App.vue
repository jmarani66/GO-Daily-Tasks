<template>
  <div id="app" class="container">
    <h1 class="title">Daily Task Planner</h1>
    <div class="task-input">
      <input v-model="newTask" @keyup.enter="addTask" placeholder="Add a new task" class="input">
      <button @click="addTask" class="button is-primary">Add</button>
    </div>
    <div v-for="task in tasks" :key="task.id" class="task">
      <div class="task-details">
        <input type="checkbox" v-model="task.done" @change="updateTask(task)">
        <span :class="{ 'done': task.done }">{{ task.title }}</span>
      </div>
      <button @click="deleteTask(task)" class="button is-danger is-small">Delete</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      tasks: [],
      newTask: ''
    };
  },
  created() {
    this.fetchTasks();
  },
  methods: {
    fetchTasks() {
      axios.get('/api/tasks')
        .then(response => {
          this.tasks = response.data;
        })
        .catch(error => {
          console.error(error);
        });
    },
    addTask() {
      if (this.newTask.trim() === '') {
        return;
      }
      axios.post('/api/tasks', { title: this.newTask, done: false })
        .then(response => {
          this.tasks.push(response.data);
          this.newTask = '';
        })
        .catch(error => {
          console.error(error);
        });
    },
    updateTask(task) {
      axios.put(`/api/tasks/${task.id}`, task)
        .catch(error => {
          console.error(error);
        });
    },
    deleteTask(task) {
      axios.delete(`/api/tasks/${task.id}`)
        .then(() => {
          this.tasks = this.tasks.filter(t => t.id !== task.id);
        })
        .catch(error => {
          console.error(error);
        });
    }
  }
};
</script>

<style>
#app {
  margin-top: 2rem;
}
.task-input {
  display: flex;
  margin-bottom: 1rem;
}
.input {
  margin-right: 1rem;
}
.task {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}
.task-details {
  display: flex;
  align-items: center;
}
.task-details input {
  margin-right: 0.5rem;
}
.done {
  text-decoration: line-through;
}
</style>
