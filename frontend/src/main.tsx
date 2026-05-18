import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ApolloClient, InMemoryCache, ApolloProvider } from "@apollo/client";
import App from "./App";
import PrivacyPolicy from "./pages/PrivacyPolicy";
import Insights from "./pages/Insights";
import "./index.css";

const apiBase = import.meta.env.VITE_API_URL ?? "";

const client = new ApolloClient({
  uri: `${apiBase}/graphql`,
  cache: new InMemoryCache(),
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Insights />} />
          <Route path="/map" element={<App />} />
          <Route path="/privacy" element={<PrivacyPolicy />} />
          <Route path="/insights" element={<Insights />} />
        </Routes>
      </BrowserRouter>
    </ApolloProvider>
  </React.StrictMode>
);
