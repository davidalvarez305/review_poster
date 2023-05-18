import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Paragraph, Sentence, Template } from "../types/general";
import { UserContext } from "../context/UserContext";
import { createUpdateSentencesFactory } from "../utils/createUpdateSentencesFactory";
import { getId } from "../utils/getId";

export default function useSentencesController() {
  const { user } = useContext(UserContext);
  const [sentences, setSentences] = useState<Sentence[]>([]);
  const [editSingleSentence, setEditSingleSentence] = useState<Sentence | null>(
    null
  );
  const { isLoading, makeRequest, error } = useFetch();
  const paragraph = useMemo((): string | undefined => {
    return window.location.pathname.split("/paragraph/")[1]
  }, []);
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

  const getSentencesByParagraph = useCallback(() => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
        url: USER_ROUTE + `/${user.id}/sentence?paragraph=${paragraph}`,
      },
      (res) => setSentences(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS, paragraph, user.id]);

  const deleteSentence = useCallback(
    (id: number) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/sentence/${[id]}?sentences=${[id]}&paragraph=${paragraph}`,
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

  const updateSentence = useCallback(
    (opts: { input: string }) => {
      const sentence = opts.input.split("\n")[0];
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/sentence/${editSingleSentence?.id}`,
          method: "PUT",
          data: {
            ...editSingleSentence,
            sentence,
          },
        },
        (res) => {
          const response: Sentence[] = res.data.data;
          setSentences(response);
          setEditSingleSentence(null);
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, editSingleSentence]
  );

  const createSentences = useCallback(
    (
      opts: { paragraph: number; template: number; sentence: string },
    ) => {
      const sentenceBody = opts.sentence.split("\n").map((sentence) => {
        return {
          paragraph_id: opts.paragraph,
          template_id: opts.template,
          sentence,
        };
      });

      if (!opts.paragraph || !opts.template) {
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
    if (paragraph) {
      getSentencesByParagraph();
    } else {
      getSentences();
    }
  }, [getSentences, getSentencesByParagraph, paragraph]);

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
    createSentences,
    updateSentence,
  };
}
