import { ChakraProvider } from "@chakra-ui/react";
import * as React from "react";
import * as ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { App } from "./App";
import { UserProvider } from "./context/UserContext";
import { ChangePassword } from "./screens/ChangePassword";
import { CreateParagraph } from "./screens/CreateParagraph";
import { CreateSentence } from "./screens/CreateSentence";
import { Dictionary } from "./screens/Dictionary";
import { Token } from "./screens/Token";
import Generate from "./screens/Generate";
import Login from "./screens/Login";
import { ParagraphsList } from "./screens/ParagraphsList";
import Register from "./screens/Register";
import { SentencesList } from "./screens/SentencesList";
import { SynonymsList } from "./screens/SynonymsList";
import User from "./screens/User";

const container = document.getElementById("root");
if (!container) throw new Error("Failed to find the root element");
const root = ReactDOM.createRoot(container);

root.render(
  <React.StrictMode>
    <ChakraProvider>
      <UserProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<App />} />
            <Route path="/dictionary" element={<Dictionary />} />
            <Route path="/create-sentence" element={<CreateSentence />} />
            <Route path="/create-paragraph" element={<CreateParagraph />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/word/:word" element={<SynonymsList />} />
            <Route path="/paragraph/:paragraph" element={<SentencesList />} />
            <Route path="/template/:template" element={<ParagraphsList />} />
            <Route path="/generate" element={<Generate />} />
            <Route path="/user" element={<User />} />
            <Route path="/change-password" element={<ChangePassword />} />
            <Route path="/token/:code" element={<Token />} />
          </Routes>
        </BrowserRouter>
      </UserProvider>
    </ChakraProvider>
  </React.StrictMode>
);
