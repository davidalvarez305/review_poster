import { Stack, Button } from "@chakra-ui/react";
import { Form } from "formik";
import SimpleInputField from "../components/SimpleInputField";
import RequestErrorMessage from "../components/RequestErrorMessage";

interface Props {
  isLoading: boolean;
  loginError: { message: string };
}

const LoginForm: React.FC<Props> = ({ isLoading, loginError }) => {
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
            Submit
          </Button>
          <RequestErrorMessage {...loginError} />
        </Stack>
      </Stack>
    </Form>
  );
};

export default LoginForm;
