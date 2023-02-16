import {
  FormControl,
  FormErrorMessage,
  FormLabel,
} from "@chakra-ui/form-control";
import { Input } from "@chakra-ui/input";
import { Box } from "@chakra-ui/layout";
import { useField } from "formik";
import React, { InputHTMLAttributes } from "react";

type Props = InputHTMLAttributes<HTMLInputElement> & {
  label: string;
  name: string;
};

const SimpleInputField: React.FC<Props> = ({
  label,
  name,
  size: _,
  ...props
}) => {
  const [field, meta] = useField(name);
  return (
    <Box>
      <FormControl>
        <FormLabel htmlFor={field.name}>{label}</FormLabel>
        <Input {...props} {...field} id={field.name} />
        {meta.error && meta.touched && (
          <FormErrorMessage>{meta.error}</FormErrorMessage>
        )}
      </FormControl>
    </Box>
  );
};

export default SimpleInputField;
