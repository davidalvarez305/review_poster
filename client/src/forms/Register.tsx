import { Box, Stack, Button } from "@chakra-ui/react";
import { Form } from "formik";
import SimpleInputField from "../components/SimpleInputField";
import RequestErrorMessage from "../components/RequestErrorMessage";

interface Props {
  isLoading: boolean;
  registerError: { message: string };
}

const RegisterForm: React.FC<Props> = ({ isLoading, registerError }) => {
  return (
    <Form>
      <Stack spacing={4}>
        <SimpleInputField
          label={"Username"}
          name="username"
          placeholder={"Enter username..."}
        />
        <SimpleInputField
          name="password"
          type="password"
          label={"Password"}
          placeholder={"Enter password..."}
        />
        <SimpleInputField
          name="email"
          type="email"
          label={"Email"}
          placeholder={"Enter email..."}
        />
        <Stack spacing={10}>
          <Button
            isLoading={isLoading}
            loadingText="Submitting"
            bg={"blue.400"}
            color={"white"}
            _hover={{
              bg: "blue.500",
            }}
            variant="outline"
            type={"submit"}
          >
            Register
          </Button>
          <RequestErrorMessage {...registerError} />
        </Stack>
      </Stack>
    </Form>
  );
};

export default RegisterForm;
