import { useState, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Sentence } from "../types/general";
import { UserContext } from "../context/UserContext";
import { createUpdateSentencesFactory } from "../utils/createUpdateSentencesFactory";

export default function useSentencesController() {
  const { user } = useContext(UserContext);
  const [sentences, setSentences] = useState<Sentence[]>([]);
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
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, sentences]
  );

  const updateSentence = useCallback(
    (sentence: Sentence) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/sentence/${sentence.id}?paragraph=${paragraph}`,
          method: "PUT",
          data: sentence,
        },
        (res) => {
          setSentences(res.data.data);
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, paragraph]
  );

  const getSentences = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setSentences(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const getSentencesByParagraph = useCallback((paragraph: string) => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
        url: USER_ROUTE + `/${user.id}/sentence?paragraph=${paragraph}`,
      },
      (res) => setSentences(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS, user.id]);

  const getSentencesByTemplate = useCallback((template: string) => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
        url: USER_ROUTE + `/${user.id}/sentence?template=${template}`,
      },
      (res) => setSentences(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS, user.id]);

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
        return { sentence, paragraph_id: sentences[0].paragraph_id };
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

  return {
    updateSentences,
    getSentences,
    setSentences,
    sentences,
    isLoading,
    error,
    deleteSentence,
    bulkUpdateSentences,
    paragraph,
    createSentences,
    updateSentence,
    getSentencesByTemplate,
    getSentencesByParagraph
  };
}
