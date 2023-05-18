import { Box, Button } from "@chakra-ui/react";
import React from "react";
import LargeInputBox from "../components/LargeInputBox";
import useLoginRequired from "../hooks/useLoginRequired";
import { centeredDiv } from "../utils/centeredDiv";
import { Formik, Form } from "formik";
import useFetch from "../hooks/useFetch";
import FormSelectComponent from "../components/FormSelectComponent";
import Layout from "../layout/Layout";
import { BottomNavigation } from "../components/BottomNavigation";
import RequestErrorMessage from "../components/RequestErrorMessage";
import { API_ROUTE } from "../constants";
import useWordsController from "../hooks/useWordsController";
import createTagFactory from "../utils/createTagFactory";

interface Props {}

export const Dictionary: React.FC<Props> = () => {
  const { createWords, words, isLoading, error } = useWordsController();
  const { makeRequest, isLoading: Loading } = useFetch();
  useLoginRequired();

  function handlePullFromChatGPT(
    values: {
      word: string;
      synonyms: string;
      id: null;
    },
    setFieldValue: (
      field: string,
      value: any,
      shouldValidate?: boolean | undefined
    ) => void
  ) {
    const tag = createTagFactory(values.word);
    makeRequest(
      {
        url: API_ROUTE + `/ai/tags?tag=${encodeURIComponent(tag)}`,
      },
      async (res) => {
        setFieldValue("synonyms", res.data.data.join("\n"));
      }
    );
  }

  return (
    <Layout>
      <Formik
        initialValues={{ word: "", synonyms: "", id: null }}
        onSubmit={(values, actions) => {
          createWords(values);
          actions.resetForm({
            values: {
              word: "",
              synonyms: "",
              id: null,
            },
          });
        }}
      >
        {({ values, setFieldValue }) => (
          <Form>
            <Box sx={{ ...centeredDiv, gap: 2, height: "100%", my: 5 }}>
              <Box sx={{ ...centeredDiv, width: "25%", height: "20%" }}>
                <FormSelectComponent
                  options={words.map((word) => word.name)}
                  name={"word"}
                />
              </Box>
              <Box>
                <LargeInputBox label="Synonyms" name="synonyms" />
              </Box>
              <Box sx={{ ...centeredDiv, flexDirection: "row", gap: 5 }}>
                <Button
                  variant={"outline"}
                  colorScheme={"teal"}
                  size={"md"}
                  type={"submit"}
                  isLoading={isLoading}
                  loadingText={"Submitting"}
                >
                  Submit
                </Button>
                {values.word.length > 0 && (
                  <Button
                    variant={"outline"}
                    colorScheme={"red"}
                    size={"md"}
                    type={"button"}
                    isLoading={Loading}
                    loadingText={"Pulling"}
                    onClick={() => handlePullFromChatGPT(values, setFieldValue)}
                  >
                    Pull From ChatGPT
                  </Button>
                )}
              </Box>
            </Box>
          </Form>
        )}
      </Formik>
      <BottomNavigation message={"Enter A Word"} path={"word"} />
      <RequestErrorMessage {...error} />
    </Layout>
  );
};
