import {
  FormControl,
  FormErrorMessage,
  FormLabel,
} from "@chakra-ui/form-control";
import { Box } from "@chakra-ui/layout";
import { Textarea } from "@chakra-ui/textarea";
import { useField } from "formik";
import React, { InputHTMLAttributes } from "react";

type Props = InputHTMLAttributes<HTMLTextAreaElement> & {
  label: string;
  name: string;
};

const LargeInputBox: React.FC<Props> = ({ label, name, size: _, ...props }) => {
  const [field, meta] = useField(name);
  return (
    <Box>
      <FormControl>
        <FormLabel
          htmlFor={field.name}
          sx={{ display: "flex", justifyContent: "center" }}
        >
          {label}
        </FormLabel>
        <Textarea {...props} {...field} minH={'450px'} minW={"600px"} />
        {meta.error && meta.touched && (
          <FormErrorMessage>{meta.error}</FormErrorMessage>
        )}
      </FormControl>
    </Box>
  );
};

export default LargeInputBox;
