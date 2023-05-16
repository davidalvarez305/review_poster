import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Word, WordFormInput } from "../types/general";
import { UserContext } from "../context/UserContext";
import createWordsFactory from "../utils/createWordsFactory";

export default function useWordsController() {
  const { user } = useContext(UserContext);
  const [words, setWords] = useState<Word[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/word`,
      method: "POST",
    };
  }, [user.id]);

  useEffect(() => {
    makeRequest({ url: USER_ROUTE + `/${user.id}/word` }, (res) =>
      setWords(res.data.data)
    );
  }, [makeRequest, setWords, user.id]);

  const updateWords = useCallback(
    (opts: WordFormInput) => {
      const wordsToCreate = createWordsFactory(opts, user.id, words);
      makeRequest({ ...FETCH_PARAMS, method: "PUT", data: wordsToCreate }, (res) =>
        setWords(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, words]
  );

  const getWords = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setWords(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const createWords = useCallback(
    (opts: WordFormInput) => {
      const wordsToCreate = createWordsFactory(opts, user.id, words);
      makeRequest({ ...FETCH_PARAMS, data: wordsToCreate }, (res) =>
        setWords(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS, words]
  );

  useEffect(() => {
    getWords();
  }, [getWords]);

  return {
    updateWords,
    getWords,
    setWords,
    createWords,
    words,
    isLoading,
    error
  };
}
