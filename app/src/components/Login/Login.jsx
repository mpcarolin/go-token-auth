import styled from "styled-components";
import { useState } from "react";

const LoginContainer = styled.div`
    display: flex;
    flex-direction: column;
`

const LoginInput = styled.input`
    width: 200px;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
    margin-bottom: 10px;
`

const LoginButton = styled.button`
    width: 200px;
    padding: 10px;
`

const login = async (username, password) => {
    const response = await fetch("http://localhost:8000/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
    });
    const data = await response.json();
    return data;
}

export const Login = ({ }) => {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const handleLogin = async () => {
        const data = await login(username, password);
        console.log(data);
    }

    return (
        <LoginContainer>
            <p>Enter your username and password to login</p>
            <LoginInput type="text" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
            <LoginInput type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
            <LoginButton onClick={handleLogin}>Login</LoginButton>
        </LoginContainer>
    )
} 