import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Paragraph, Sentence, Template } from "../types/general";
import { UserContext } from "../context/UserContext";
import { useLocation } from "react-router-dom";
import { createUpdateSentencesFactory } from "../utils/createUpdateSentencesFactory";
import { getId } from "../utils/getId";

export default function useSentencesController() {
  const { user } = useContext(UserContext);
  const [sentences, setSentences] = useState<Sentence[]>([]);
  const [editSingleSentence, setEditSingleSentence] = useState<Sentence | null>(
    null
  );
  const { isLoading, makeRequest, error } = useFetch();
  const location = useLocation();
  const paragraph = useMemo(
    () => location.pathname.split("/paragraph/")[1],
    [location.pathname]
  );
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/sentence?paragraph=${paragraph}`,
      method: "POST",
    };
  }, [user.id, paragraph]);

  const updateSentences = useCallback(
    (opts: { input: string }, paragraph_id: number, paragraph: string) => {
      const sentencesToUpdate = createUpdateSentencesFactory(
        sentences,
        opts.input.split("\n"),
        paragraph_id
      );
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/sentence?paragraph=${paragraph}`,
          method: "PUT",
          data: sentencesToUpdate,
        },
        (res) => {
          setSentences(res.data.data);
          setEditSingleSentence(null);
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, sentences]
  );

  const getSentences = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setSentences(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const deleteSentence = useCallback(
    (id: number) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/sentence?paragraph=${paragraph}&sentences=${[id]}`,
          method: "DELETE",
        },
        (res) => setSentences(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS, paragraph]
  );

  const bulkUpdateSentences = useCallback(
    (values: { input: string }) => {
      let body = values.input.split("\n").map((sentence) => {
        return { sentence, paragraph_id: sentences[0].id };
      });
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/sentence/bulk?paragraph=${paragraph}`,
          method: "POST",
          data: body,
        },
        (res) => setSentences(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id, paragraph, sentences]
  );

  const createSentences = useCallback(
    (
      opts: { paragraph: string; template: string; sentence: string },
      paragraphs: Paragraph[],
      templates: Template[]
    ) => {
      const paragraph_id = getId(opts.paragraph, paragraphs, "name");
      const template_id = getId(opts.template, templates, "name");
      const sentenceBody = opts.sentence.split("\n").map((sentence) => {
        return {
          paragraph_id,
          template_id,
          sentence,
        };
      });

      if (!paragraph_id || !template_id) {
        return;
      }

      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/sentence`,
          method: "POST",
          data: sentenceBody,
        },
        (res) => setSentences(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id]
  );

  useEffect(() => {
    getSentences();
  }, [getSentences]);

  return {
    updateSentences,
    getSentences,
    setSentences,
    sentences,
    isLoading,
    error,
    setEditSingleSentence,
    editSingleSentence,
    deleteSentence,
    bulkUpdateSentences,
    paragraph,
    createSentences
  };
}
