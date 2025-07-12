document.addEventListener('DOMContentLoaded', () => {
  const newTaskInput = document.getElementById('new-task');
  const addTaskButton = document.getElementById('add-task');
  const taskList = document.getElementById('task-list');

  const apiUrl = '/api/tasks';

  const fetchTasks = async () => {
    try {
      const response = await fetch(apiUrl);
      const tasks = await response.json();
      renderTasks(tasks);
    } catch (error) {
      console.error('Error fetching tasks:', error);
    }
  };

  const renderTasks = (tasks) => {
    taskList.innerHTML = '';
    tasks.forEach(task => {
      const taskItem = document.createElement('li');
      taskItem.innerHTML = `
        <div class="task-details">
          <input type="checkbox" ${task.done ? 'checked' : ''} data-id="${task.id}">
          <span class="${task.done ? 'done' : ''}">${task.title}</span>
        </div>
        <button class="delete-btn" data-id="${task.id}">Delete</button>
      `;
      taskList.appendChild(taskItem);
    });
  };

  const addTask = async () => {
    const title = newTaskInput.value.trim();
    if (title === '') {
      return;
    }

    try {
      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title, done: false }),
      });
      const newTask = await response.json();
      newTaskInput.value = '';
      fetchTasks();
    } catch (error) {
      console.error('Error adding task:', error);
    }
  };

  const updateTask = async (id, done, title) => {
    try {
      await fetch(`${apiUrl}/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id: parseInt(id), done, title }),
      });
      fetchTasks();
    } catch (error) {
      console.error('Error updating task:', error);
    }
  };

  const deleteTask = async (id) => {
    try {
      await fetch(`${apiUrl}/${id}`, {
        method: 'DELETE',
      });
      fetchTasks();
    } catch (error) {
      console.error('Error deleting task:', error);
    }
  };

  addTaskButton.addEventListener('click', addTask);
  newTaskInput.addEventListener('keyup', (event) => {
    if (event.key === 'Enter') {
      addTask();
    }
  });

  taskList.addEventListener('click', (event) => {
    if (event.target.matches('.delete-btn')) {
      const id = event.target.dataset.id;
      deleteTask(id);
    } else if (event.target.matches('input[type="checkbox"]')) {
      const id = event.target.dataset.id;
      const done = event.target.checked;
      const title = event.target.nextElementSibling.textContent;
      updateTask(id, done, title);
    }
  });

  fetchTasks();
});
