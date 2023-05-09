import {
  Flex,
  Box,
  Stack,
  Heading,
  Text,
  useColorModeValue,
} from "@chakra-ui/react";
import useFetch from "../hooks/useFetch";
import LoginForm from "../forms/Login";
import { LOGIN_ROUTE } from "../constants";
import LoginOrRegister from "../components/LoginOrRegister";
import { Formik } from "formik";
import { useNavigate } from "react-router";
import { UserContext } from "../context/UserContext";
import { useContext } from "react";
import useLoginRequired from "../hooks/useLoginRequired";

const LoginScreen = () => {
  useLoginRequired();
  const { SetUser } = useContext(UserContext);
  const { isLoading, makeRequest, error } = useFetch();
  const navigate = useNavigate();

  function handleRequest(values: { username: string; password: string }) {
    if (values.username === "" || values.password === "") {
      return;
    }
    makeRequest(
      {
        url: LOGIN_ROUTE,
        method: "POST",
        data: values,
      },
      (res) => {
        // Save the user in context after successful login.
        SetUser(res.data.data);
        navigate("/");
      }
    );
  }

  return (
    <Flex minH={"100vh"} align={"top"} justify={"center"}>
      <Stack spacing={8} mx={"auto"} w={"25%"} py={12} px={6}>
        <Stack align={"center"}>
          <Heading fontSize={"4xl"}>Login</Heading>
          <Text fontSize={"lg"} color={"gray.600"}>
            To the moon 🚀😎💻
          </Text>
        </Stack>
        <Box
          rounded={"lg"}
          bg={useColorModeValue("white", "gray.700")}
          boxShadow={"lg"}
          p={8}
        >
          <Formik
            initialValues={{ username: "", password: "" }}
            onSubmit={handleRequest}
          >
            <LoginForm isLoading={isLoading} loginError={error} />
          </Formik>
          <LoginOrRegister text={"Create Account"} navigatePage={"register"} />
        </Box>
      </Stack>
    </Flex>
  );
};
export default LoginScreen;
