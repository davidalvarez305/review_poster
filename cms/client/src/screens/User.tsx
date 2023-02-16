import { Box, Button, IconButton } from "@chakra-ui/react";
import { Form, Formik } from "formik";
import React, { useContext } from "react";
import { useNavigate } from "react-router";
import SimpleInputField from "../components/SimpleInputField";
import { USER_ROUTE } from "../constants";
import { UserContext } from "../context/UserContext";
import { emptyUser } from "../hooks/useAuth";
import useFetch from "../hooks/useFetch";
import useLoginRequired from "../hooks/useLoginRequired";
import Layout from "../layout/Layout";
import { User as UserType } from "../types/general";
import { FiRefreshCcw } from "react-icons/fi";

const User: React.FC = () => {
  const { user, SetUser, Logout } = useContext(UserContext);
  const { makeRequest, isLoading } = useFetch();
  const navigate = useNavigate();
  useLoginRequired();

  function handleDeleteAccount() {
    makeRequest(
      {
        url: USER_ROUTE,
        method: "DELETE",
      },
      (_) => {
        SetUser(emptyUser);
      }
    );
  }

  function handleRequest(values: UserType) {
    makeRequest(
      {
        url: USER_ROUTE,
        method: "PUT",
        data: values,
      },
      (res) => {
        SetUser(res.data.data);
      }
    );
  }

  return (
    <Layout>
      <Box sx={{ my: 10 }}>
        <Formik initialValues={user} onSubmit={handleRequest}>
          {({ values }) => (
            <Form>
              {Object.keys(user).map((u) => (
                <div
                  key={String(u)}
                  style={{
                    marginBottom: 15,
                    display: "flex",
                    flexDirection: "row",
                    justifyContent: "center",
                    alignItems: "center",
                    paddingLeft: 40,
                  }}
                >
                  <SimpleInputField
                    label={String(u)}
                    name={`${u}`}
                    width={250}
                  />
                  {u === "token" ? (
                    <IconButton
                      icon={<FiRefreshCcw />}
                      aria-label={"refresh"}
                      style={{ marginTop: 30, marginLeft: 10 }}
                      onClick={() => handleRequest(values)}
                      isLoading={isLoading}
                    />
                  ) : undefined}
                </div>
              ))}
              <Button
                variant={"outline"}
                colorScheme={"green"}
                sx={{ my: 5 }}
                type={"submit"}
                isLoading={isLoading}
                loadingText={"Submitting"}
              >
                Save Changes
              </Button>
              <Button
                variant={"outline"}
                colorScheme={"orange"}
                sx={{ my: 5, ml: 5 }}
                onClick={() => {
                  navigate("/change-password");
                }}
              >
                Change Password
              </Button>
            </Form>
          )}
        </Formik>
      </Box>
      <Box
        display={"flex"}
        flexDir={"row"}
        gap={5}
        justifyContent={"center"}
        alignItems={"center"}
      >
        <Button
          variant={"outline"}
          colorScheme={"blue"}
          sx={{ my: 5 }}
          isLoading={isLoading}
          loadingText={"Submitting"}
          onClick={() => Logout()}
        >
          Logout
        </Button>
        <Button
          variant={"outline"}
          colorScheme={"red"}
          sx={{ my: 5 }}
          isLoading={isLoading}
          loadingText={"Submitting"}
          onClick={() => handleDeleteAccount()}
        >
          Delete Account
        </Button>
      </Box>
    </Layout>
  );
};

export default User;
