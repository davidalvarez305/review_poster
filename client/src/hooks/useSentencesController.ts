import { useState, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Paragraph, Sentence, Template } from "../types/general";
import { UserContext } from "../context/UserContext";
import createUpdateSentencesFactory from "../utils/createUpdateSentencesFactory";
import deleteUserParagraphSentencesByTemplateFactory from "../utils/deleteUserParagraphSentencesByTemplateFactory";
import createSentencesFactory from "../utils/createSentencesFactory";

export default function useSentencesController() {
  const { user } = useContext(UserContext);
  const [sentences, setSentences] = useState<Sentence[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/template/`,
      method: "POST",
    };
  }, [user.id]);

  const deleteUserParagraphSentencesByTemplate = useCallback(
    (values: { input: string }, template: string, paragraph: string) => {
      const body = deleteUserParagraphSentencesByTemplateFactory(
        values,
        sentences
      );
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/template/${template}/paragraph/${paragraph}/sentence`,
          method: "DELETE",
          data: body,
        },
        (res) => setSentences(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id, sentences]
  );

  const updateUserParagraphSentencesByTemplate = useCallback(
    (values: { input: string }, paragraph: string, template: string) => {
      const sentencesToUpdate = createUpdateSentencesFactory(values, sentences);
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/template/${template}/paragraph/${paragraph}/sentence`,
          method: "PATCH",
          data: sentencesToUpdate,
        },
        (res) => {
          setSentences(res.data.data);
          deleteUserParagraphSentencesByTemplate(
            { ...values },
            paragraph,
            template
          );
        }
      );
    },
    [
      makeRequest,
      FETCH_PARAMS,
      user.id,
      sentences,
      deleteUserParagraphSentencesByTemplate,
    ]
  );

  const updateParagraphSentenceByTemplate = useCallback(
    (sentence: Sentence, template: string, paragraph: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/template/${template}/paragraph/${paragraph}/sentence/${sentence.id}`,
          method: "PUT",
          data: sentence,
        },
        (res) => {
          setSentences(res.data.data);
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id]
  );

  const deleteUserParagraphSentenceByTemplate = useCallback(
    (id: number, paragraph: string, template: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/template/${template}/paragraph/${paragraph}/sentence/${id}`,
          method: "DELETE",
        },
        (res) => setSentences(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS]
  );

  const createUserParagraphSentencesByTemplate = useCallback(
    (values: { input: string }, paragraph: Paragraph, template: Template) => {
      const body = createSentencesFactory(values, paragraph);
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/template/${template.name}/paragraph/${paragraph.name}/sentence`,
          method: "POST",
          data: body,
        },
        (res) => setSentences(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id]
  );

  const getUserParagraphSentencesByTemplate = useCallback(
    (paragraph: string, template: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          method: "GET",
          url: USER_ROUTE + `/${user.id}/template/${template}/paragraph/${paragraph}/sentence`,
        },
        (res) => setSentences(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id]
  );

  const getUserSentencesByTemplate = useCallback((template: string) => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
        url: USER_ROUTE + `/${user.id}/template/${template}/sentence`,
      },
      (res) => setSentences(res.data.data)
    );
  }, [FETCH_PARAMS, makeRequest, user.id]);

  return {
    setSentences,
    sentences,
    isLoading,
    error,
    updateUserParagraphSentencesByTemplate,
    deleteUserParagraphSentenceByTemplate,
    deleteUserParagraphSentencesByTemplate,
    createUserParagraphSentencesByTemplate,
    updateParagraphSentenceByTemplate,
    getUserParagraphSentencesByTemplate,
    getUserSentencesByTemplate
  };
}
