import React, { useState } from 'react';
import axios from 'axios';
import { useAuth } from './context/AuthContext';

const Login = () => {
    const { setAuth } = useAuth(); // получение функции для установки авторизации
    const [username, setUsername] = useState(''); 
    const [password, setPassword] = useState(''); 
    const [error, setError] = useState(''); 

    const handleLogin = async () => {
        try {
            const response = await axios.post('http://localhost:8080/login', { username, password }); // Отправляем запрос на сервер
            setAuth({ token: response.data.token }); // установка токена авторизации
            setError(''); // очистка ошибки при успешном входе
        } catch (error) {
            console.error("Login error", error); // вывод ошибки в консоль
            setError('Ошибка при входе. Проверьте логин и пароль.'); // установка сообщения об ошибке
        }
    };

    return (
        <div>
            <h2>Login</h2>
            <input
                type="text"
                value={username} // установка значение имени пользователя
                onChange={(e) => setUsername(e.target.value)} // обновление имени пользователя
                placeholder="Username" //  для поля ввода
            />
            <input
                type="password"
                value={password} // установка значение пароля
                onChange={(e) => setPassword(e.target.value)} // обновление пароля
                placeholder="Password" //  для поля ввода
            />
            <button onClick={handleLogin}>Login</button> {/* кнопка  входа */}
            {error && <div style={{ color: 'red' }}>{error}</div>} {/*  сообщение об ошибке, если оно есть */}
        </div>
    );
};

export default Login; 
