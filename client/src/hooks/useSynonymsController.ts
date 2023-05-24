import { useState, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Synonym, Word } from "../types/general";
import { UserContext } from "../context/UserContext";
import updateUserSynonymsByWordFactory from "../utils/updateUserSynonymsByWordFactory";
import createSynonymsFactory from "../utils/createSynonymsFactory";

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

  const updateUserSynonymByWord = useCallback(
    (synonym: Synonym, word: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/word/${word}/synonym/${synonym.id}`,
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

  const getUserSynonymsByWord = useCallback(
    (wordString: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          method: "GET",
          url: USER_ROUTE + `/${user.id}/word/${wordString}/synonym`,
        },
        (res) => setSynonyms(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id]
  );

  const deleteUserSynonymByWord = useCallback(
    (id: number, word: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/word/${word}/synonym/${id}`,
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

  const createUserSynonymsByWord = useCallback(
    (values: { input: string }, word: Word) => {
      const body = createSynonymsFactory(values, word);
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/word/${word.name}/synonym`,
          method: "POST",
          data: body,
        },
        (res) => setSynonyms(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id]
  );

  return {
    setSynonyms,
    synonyms,
    isLoading,
    error,
    deleteUserSynonymByWord,
    updateUserSynonymsByWord,
    updateUserSynonymByWord,
    getUserSynonymsByWord,
    createUserSynonymsByWord,
  };
}
