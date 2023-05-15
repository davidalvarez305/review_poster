import { Box, Button, useToast } from "@chakra-ui/react";
import React, { useContext, useEffect, useState } from "react";
import LargeInputBox from "../components/LargeInputBox";
import useLoginRequired from "../hooks/useLoginRequired";
import { centeredDiv } from "../utils/centeredDiv";
import { Formik, Form } from "formik";
import useFetch from "../hooks/useFetch";
import FormSelectComponent from "../components/FormSelectComponent";
import { createTag } from "../utils/createTag";
import Layout from "../layout/Layout";
import { BottomNavigation } from "../components/BottomNavigation";
import { UserContext } from "../context/UserContext";
import RequestErrorMessage from "../components/RequestErrorMessage";
import { getId } from "../utils/getId";
import { Word } from "../types/general";
import { API_ROUTE, USER_ROUTE } from "../constants";

interface Props {}

export const Dictionary: React.FC<Props> = () => {
  const { user } = useContext(UserContext);
  const { isLoading, makeRequest, error } = useFetch();
  const [words, setWords] = useState<Word[]>([]);
  const toast = useToast();
  useLoginRequired();

  function handlePullFromChatGPT(values: {
    word: string;
    synonyms: string;
    id: null;
  }, setFieldValue: (field: string, value: any, shouldValidate?: boolean | undefined) => void) {
    const tag = createTag(values.word);
    makeRequest(
      {
        url: API_ROUTE + `/ai/tags?tag=${encodeURIComponent(tag)}`,
      },
      async (res) => {
        setFieldValue("synonyms", res.data.data.join("\n"));
      }
    );
  }

  function handleSubmit(values: {
    word: string;
    synonyms: string;
    id: number | null;
  }) {
    const wordStruct = {
      id: getId(values.word, words, "name"),
      word: values.word,
      tag: createTag(values.word),
      user_id: user.id,
      synonyms: values.synonyms.split("\n"),
    };
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/word`,
        method: "POST",
        data: wordStruct,
      },
      async (res) => {
        toast({
          title: "Success!",
          description: "Word has been submitted",
          status: "success",
          isClosable: true,
          duration: 5000,
          variant: "left-accent",
        });
        setWords(res.data.data);
      }
    );
  }

  useEffect(() => {
    makeRequest(
      {
        url: USER_ROUTE + `/${user.id}/word`,
      },
      (res) => {
        setWords(res.data.data);
      }
    );
  }, [makeRequest, user.id]);

  return (
    <Layout>
      <Formik
        initialValues={{ word: "", synonyms: "", id: null }}
        onSubmit={(values, actions) => {
          handleSubmit(values);
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
                    isLoading={isLoading}
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
