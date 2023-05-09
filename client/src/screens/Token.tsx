import { Box, Button } from "@chakra-ui/react";
import { Form, Formik } from "formik";
import React, { useContext } from "react";
import { useParams } from "react-router-dom";
import SimpleInputField from "../components/SimpleInputField";
import { USER_ROUTE } from "../constants";
import { UserContext } from "../context/UserContext";
import useFetch from "../hooks/useFetch";
import useLoginRequired from "../hooks/useLoginRequired";
import Layout from "../layout/Layout";
import RequestErrorMessage from "../components/RequestErrorMessage";

export const Token = () => {
  const { code } = useParams();
  const { Logout } = useContext(UserContext);
  const { makeRequest, isLoading, error } = useFetch();
  useLoginRequired();

  const centered = {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    height: "80vh",
    weight: "100%",
  };

  function handleChangePassword(values: { newPassword: string }) {
    makeRequest(
      {
        url: USER_ROUTE + "/change-password/" + code,
        method: "PUT",
        data: values,
      },
      (_) => {
        Logout();
      }
    );
  }

  return (
    <Layout>
      <Box sx={{ ...centered, flexDirection: "column" }}>
        <Formik
          initialValues={{ newPassword: "" }}
          onSubmit={handleChangePassword}
        >
          <Form>
            <SimpleInputField
              label={"New Password"}
              name={"newPassword"}
              type={"password"}
            />
            <Button
              variant={"outline"}
              colorScheme={"teal"}
              isLoading={isLoading}
              loadingText={"Submitting"}
              type={"submit"}
              marginY={5}
            >
              Submit
            </Button>
          </Form>
        </Formik>
        {error.message.length > 0 && <RequestErrorMessage {...error} />}
      </Box>
    </Layout>
  );
};
