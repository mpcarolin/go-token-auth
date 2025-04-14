import React from "react";
import { createRoot } from "react-dom/client";
import styled, { createGlobalStyle } from "styled-components";
import { Login } from "./components/Login/Login";

const GlobalStyle = createGlobalStyle`
  body {
    font-family: Verdana, Geneva, sans-serif;
    font-size: 14px;
    line-height: 1.6;
    color: #333;
    background-color: #f6f6ef;
  }
`;

const AppContainer = styled.div`
    display: flex;
    padding: 20px;
`;

const App = () => {
    return (
        <>
            <GlobalStyle />
            <AppContainer>
                <Login />
            </AppContainer>
        </>
    )
}

const root = createRoot(document.getElementById("root"))

root.render(<App />)