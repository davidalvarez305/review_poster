import { Box, Button } from "@chakra-ui/react";
import React, { useState } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "../hooks/useFetch";
import useLoginRequired from "../hooks/useLoginRequired";
import Layout from "../layout/Layout";

export const ChangePassword = () => {
  useLoginRequired()
  const { makeRequest, isLoading } = useFetch();
  const [msg, setMsg] = useState("");

  const centered = {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    height: "80vh",
    weight: "100%",
  };

  // Request for Generating a New Code from Backend
  function handleRequest() {
    makeRequest(
      {
        url: USER_ROUTE + "/change-password/",
        method: "GET",
      },
      (_) => {
        setMsg("Click the link in your e-mail within 5 minutes.");
      }
    );
  }

  if (msg.length > 0) {
    return (
      <Layout>
        <Box
          sx={{
            ...centered,
          }}
        >
          {msg}
        </Box>
      </Layout>
    );
  }

  return (
    <Layout>
      <Box
        sx={{
          ...centered,
        }}
      >
        <Button
          variant={"outline"}
          colorScheme={"teal"}
          isLoading={isLoading}
          loadingText={"Submitting"}
          onClick={() => handleRequest()}
        >
          Request Token
        </Button>
      </Box>
    </Layout>
  );
};
