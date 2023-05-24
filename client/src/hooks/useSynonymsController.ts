import { useState, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Synonym } from "../types/general";
import { UserContext } from "../context/UserContext";
import { createUpdateSynonymsFactory } from "../utils/createUpdateSynonymsFactory";
import updateUserSynonymsByWordFactory from "../utils/updateUserSynonymsByWordFactory";

export default function useSynonymsController() {
  const { user } = useContext(UserContext);
  const [synonyms, setSynonyms] = useState<Synonym[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/word/`,
      method: "POST",
    };
  }, [user.id]);

  const updateSynonym = useCallback(
    (synonym: Synonym, word: string) => {
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
    [makeRequest, FETCH_PARAMS, user.id]
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

  const getUserSynonymsByWord = useCallback((wordString: string) => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
        url: USER_ROUTE + `/${user.id}/word/${wordString}/synonym`,
      },
      (res) => setSynonyms(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS, user.id]);

  const deleteSynonym = useCallback(
    (id: number, word: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/synonym/${[id]}?word=${word}&synonyms=${[id]}`,
          method: "DELETE",
        },
        (res) => setSynonyms(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS]
  );

  const updateUserSynonymsByWord = useCallback(
    (values: { input: string }, word: string) => {
      const body = updateUserSynonymsByWordFactory(values, synonyms);
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/word/${word}/synonym`,
          method: "PUT",
          data: body,
        },
        (res) => setSynonyms(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id, synonyms]
  );

  return {
    updateSynonyms,
    getSynonyms,
    setSynonyms,
    synonyms,
    isLoading,
    error,
    deleteSynonym,
    updateUserSynonymsByWord,
    updateSynonym,
    getUserSynonymsByWord
  };
}
