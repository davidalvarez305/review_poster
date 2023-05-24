import {
  FormControl,
  FormLabel,
  Box,
  FormErrorMessage,
} from "@chakra-ui/react";
import { useEffect, useMemo, useState } from "react";
import CreatableSelect from "react-select/creatable";
import { capitalizeFirstLetter } from "../utils/capitalizeFirstLetter";
import { useField, useFormikContext } from "formik";

type SelectType = { value: string; label: string };

interface Props {
  options: any[];
  name: string;
}

const FormSelectComponent: React.FC<Props> = ({ options, name }) => {
  let initialProps = useMemo(
    () =>
      Object.create({
        value: "",
        label: "",
      }),
    []
  );
  const { setFieldValue, values } = useFormikContext();

  const [field, meta] = useField(name);

  const [selectedValue, setSelectedValue] = useState<null | SelectType>(
    initialProps
  );

  useEffect(() => {
    if ((values as any).word === "") {
      setSelectedValue(initialProps);
    }
  }, [values, initialProps]);

  return (
    <Box
      sx={{
        ml: 2,
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        width: 250,
      }}
    >
      <FormControl>
        <FormLabel
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
          htmlFor={field.name}
        >
          {capitalizeFirstLetter(name)}
        </FormLabel>
        <CreatableSelect
          name={field.name}
          placeholder={""}
          value={selectedValue}
          onChange={(e) => {
            setSelectedValue(e);
            if (field.name === "word") {
              setFieldValue(field.name, e?.label);
            } else {
              setFieldValue(field.name, e?.value);
            }
          }}
          options={options.map((op) => {
            return {
              value: op.id,
              label: capitalizeFirstLetter(op.name),
            };
          })}
        />
        {meta.error && meta.touched && (
          <FormErrorMessage>{meta.error}</FormErrorMessage>
        )}
      </FormControl>
    </Box>
  );
};

export default FormSelectComponent;
