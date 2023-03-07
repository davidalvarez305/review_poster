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
import { USER_ROUTE } from "../constants";

interface Props {}

export const Dictionary: React.FC<Props> = () => {
  const { user } = useContext(UserContext);
  const { isLoading, makeRequest, cancelToken, error } = useFetch();
  const [words, setWords] = useState<Word[]>([]);
  const toast = useToast();
  useLoginRequired();

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
        url: USER_ROUTE + `${user.id}/word`,
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
        url: USER_ROUTE + `${user.id}/word`,
      },
      (res) => {
        setWords(res.data.data);
      }
    );
    return () => {
      cancelToken.cancel();
    };
  }, [cancelToken, makeRequest]);

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
          </Box>
        </Form>
      </Formik>
      <BottomNavigation message={"Enter A Word"} path={"word"} />
      <RequestErrorMessage {...error} />
    </Layout>
  );
};
