import {
  Flex,
  Box,
  Stack,
  Heading,
  Text,
  useColorModeValue,
} from "@chakra-ui/react";
import useFetch from "../hooks/useFetch";
import { REGISTER_ROUTE } from "../constants";
import RegisterForm from "../forms/Register";
import LoginOrRegister from "../components/LoginOrRegister";
import { Formik } from "formik";
import { useNavigate } from "react-router";

export default function Register() {
  const { isLoading, makeRequest, error } = useFetch();
  const navigate = useNavigate();

  function handleRequest(values: {
    username: string;
    password: string;
    email: string;
  }) {
    if (
      values.username === "" ||
      values.password === "" ||
      values.email === ""
    ) {
      return;
    }
    makeRequest(
      {
        url: REGISTER_ROUTE,
        method: "POST",
        data: values,
      },
      async (res) => {
        if (res.data.data) {
          navigate("../login");
        }
      }
    );
  }

  return (
    <Flex minH={"100vh"} align={"top"} justify={"center"}>
      <Stack spacing={8} mx={"auto"} w={"25%"} py={12} px={6}>
        <Stack align={"center"}>
          <Heading fontSize={"4xl"}>Create An Account</Heading>
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
            initialValues={{
              username: "",
              password: "",
              email: "",
            }}
            onSubmit={handleRequest}
          >
            <RegisterForm isLoading={isLoading} registerError={error} />
          </Formik>
          <LoginOrRegister
            text={"Have an account? Login."}
            navigatePage={"login"}
          />
        </Box>
      </Stack>
    </Flex>
  );
}
