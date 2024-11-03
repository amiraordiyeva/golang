import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import { useAuth } from './context/AuthContext';
import Login from './Login';

const ItemList = () => {
    const { auth } = useAuth();
    const [tasks, setTasks] = useState([]); // хранит список задач
    const [taskName, setTaskName] = useState(''); // имя новой задачи
    const [loading, setLoading] = useState(true); // статус загрузки
    const [editingTask, setEditingTask] = useState(null); // таск который мы редактируем

    const fetchTasks = useCallback(async () => {
        setLoading(true); // установка загрузки в правду
        try {
            const response = await axios.get('http://localhost:8080/tasks', {
                headers: { Authorization: `Bearer ${auth.token}` } // добавление токена авторизации
            });
            setTasks(response.data); // сохранение задачи из ответа
        } catch (error) {
            console.error("Error fetching tasks", error); // вывод ошибки 
        } finally {
            setLoading(false); // загрузка закончена
        }
    }, [auth.token]);

    useEffect(() => {
        if (auth.token) {
            fetchTasks(); // если есть токен загружаем задачи
        }
    }, [auth.token, fetchTasks]);

    const createTask = async () => {
        if (!taskName) return; // если имя задачи пустое ничего 

        try {
            const response = await axios.post('http://localhost:8080/tasks', { name: taskName }, {
                headers: { Authorization: `Bearer ${auth.token}` } // токен для авторизации
            });
            setTasks([...tasks, response.data]); // добавление новой задачи в список
            setTaskName(''); // очищение поля ввода
        } catch (error) {
            console.error("Error creating task", error); // ошибка при создании задачи
        }
    };

    const deleteTask = async (taskId) => {
        try {
            await axios.delete(`http://localhost:8080/tasks/${taskId}`, {
                headers: { Authorization: `Bearer ${auth.token}` } // токен для авторизации
            });
            fetchTasks(); // загрузка задачи снова после удаления
        } catch (error) {
            console.error("Error deleting task", error); // ошибка при удалении задачи
        }
    };

    const editTask = (task) => {
        setTaskName(task.name); // заполнение поля ввода именем задачи
        setEditingTask(task); // установка редактируемой задачи
    };

    const updateTask = async () => {
        if (!taskName || !editingTask) return; // если имя пустое или нет редактируемой задачи

        try {
            await axios.put(`http://localhost:8080/tasks/${editingTask.id}`, { name: taskName }, {
                headers: { Authorization: `Bearer ${auth.token}` } // токен для авторизации
            });
            setTaskName(''); // очищение поля ввода
            setEditingTask(null); // сбрасываем редактируемую задачу
            fetchTasks(); // обновление списка задач
        } catch (error) {
            console.error("Error updating task", error); // ошибка при обновлении задачи
        }
    };

    // если токен отсутствует отображение страницы входа
    if (!auth.token) {
        return <Login />;
    }

    return (
        <div>
            <h2>Task List</h2>
            {loading ? <p>Loading...</p> : ( // если загрузка показывает сообщение
                <ul>
                    {tasks.map(task => ( // отображение списка задач
                        <li key={task.id}>
                            {task.name}
                            <button onClick={() => editTask(task)}>Edit</button>
                            <button onClick={() => deleteTask(task.id)}>Delete</button>
                        </li>
                    ))}
                </ul>
            )}
            <input
                type="text"
                value={taskName} // установка текущее значения поля
                onChange={(e) => setTaskName(e.target.value)} // обновление имени задачи
                placeholder="New task" //  для поле ввода
            />
            <button onClick={editingTask ? updateTask : createTask}>
                {editingTask ? 'Update Task' : 'Add Task'} 
            </button>
        </div>
    );
};

export default ItemList;
