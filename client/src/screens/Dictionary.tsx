import { Box, Button } from "@chakra-ui/react";
import React, { useContext, useEffect } from "react";
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
import { WordFormInput } from "../types/general";
import { UserContext } from "../context/UserContext";

interface Props {}

export const Dictionary: React.FC<Props> = () => {
  const { user } = useContext(UserContext);
  const { createUserWords, words, isLoading, error, getUserWords } =
    useWordsController();
  const { makeRequest, isLoading: Loading } = useFetch();
  useLoginRequired();

  function handlePullFromChatGPT(
    values: WordFormInput,
    setFieldValue: (
      field: string,
      value: any,
      shouldValidate?: boolean | undefined
    ) => void
  ) {
    const tag = createTagFactory(values.name);
    makeRequest(
      {
        url: API_ROUTE + `/ai/tags?tag=${encodeURIComponent(tag)}`,
      },
      async (res) => {
        setFieldValue("synonyms", res.data.data.join("\n"));
      }
    );
  }

  useEffect(() => {
    user.id && getUserWords();
  }, [getUserWords, user.id]);

  return (
    <Layout>
      <Formik
        initialValues={{ name: "", synonyms: "", id: null }}
        onSubmit={(values, actions) => {
          if (values.name.length > 0 && values.synonyms.length > 0) {
            createUserWords(values);
            actions.resetForm({
              values: {
                name: "",
                synonyms: "",
                id: null,
              },
            });
          }
        }}
      >
        {({ values, setFieldValue }) => (
          <Form>
            <Box sx={{ ...centeredDiv, gap: 2, height: "100%", my: 5 }}>
              <Box sx={{ ...centeredDiv, width: "25%", height: "20%" }}>
                <FormSelectComponent options={words} name={"name"} />
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
                {values.name.length > 0 && (
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
