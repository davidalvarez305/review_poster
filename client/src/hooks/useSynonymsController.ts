import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Synonym } from "../types/general";
import { UserContext } from "../context/UserContext";
import { createUpdateSynonymsFactory } from "../utils/createUpdateSynonymsFactory";

export default function useSynonymsController() {
  const { user } = useContext(UserContext);
  const [synonyms, setSynonyms] = useState<Synonym[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const word = useMemo((): string | undefined => {
    return window.location.pathname.split("/word/")[1]
  }, []);
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/synonym?word=${word}`,
      method: "POST",
    };
  }, [user.id, word]);

  const updateSynonym = useCallback(
    (synonym: Synonym) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/synonym/${synonym.id}?word=${word}`,
          method: "PUT",
          data: synonym,
        },
        (res) => {
          setSynonyms(res.data.data);
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, word]
  );

  const updateSynonyms = useCallback(
    (opts: { input: string }, word_id: number, word: string) => {
      const synonymsToUpdate = createUpdateSynonymsFactory(
        synonyms,
        opts.input.split("\n"),
        word_id
      );
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/synonym?word=${word}`,
          method: "PUT",
          data: synonymsToUpdate,
        },
        (res) => {
          setSynonyms(res.data.data);
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, synonyms]
  );

  const getSynonyms = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setSynonyms(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const getSynonymsByWord = useCallback(() => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
        url: USER_ROUTE + `/${user.id}/synonym?word=${word}`,
      },
      (res) => setSynonyms(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS, word, user.id]);

  const deleteSynonym = useCallback(
    (id: number) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/synonym/${[id]}?word=${word}&synonyms=${[id]}`,
          method: "DELETE",
        },
        (res) => setSynonyms(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS, word]
  );

  const bulkUpdateSynonyms = useCallback(
    (values: { input: string }) => {
      let body = values.input.split("\n").map((synonym) => {
        return { synonym, word_id: synonyms[0].word_id };
      });
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/synonym/bulk?word=${word}`,
          method: "POST",
          data: body,
        },
        (res) => setSynonyms(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id, word, synonyms]
  );

  useEffect(() => {
    if (word) {
      getSynonymsByWord();
    } else {
      getSynonyms();
    }
  }, [getSynonyms, getSynonymsByWord, word]);

  return {
    updateSynonyms,
    getSynonyms,
    setSynonyms,
    synonyms,
    isLoading,
    error,
    deleteSynonym,
    bulkUpdateSynonyms,
    updateSynonym,
    word,
  };
}
