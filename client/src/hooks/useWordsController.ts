import { useState, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Word, WordFormInput } from "../types/general";
import { UserContext } from "../context/UserContext";
import createWordsFactory from "../utils/createWordsFactory";
import useSynonymsController from "./useSynonymsController";

export default function useWordsController() {
  const { user } = useContext(UserContext);
  const [words, setWords] = useState<Word[]>([]);
  const { createUserSynonymsByWord } = useSynonymsController();
  const { isLoading, makeRequest, error } = useFetch();
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/word`,
      method: "POST",
    };
  }, [user.id]);

  const updateUserWords = useCallback(
    (opts: WordFormInput) => {
      const wordsToCreate = createWordsFactory(opts, user.id, words);
      makeRequest(
        { ...FETCH_PARAMS, method: "PUT", data: wordsToCreate },
        (res) => setWords(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, words]
  );

  const getUserWords = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setWords(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const createUserWords = useCallback(
    (opts: WordFormInput) => {
      const wordsToCreate = createWordsFactory(opts, user.id, words);
      makeRequest({ ...FETCH_PARAMS, data: wordsToCreate }, (res) => {
        const word: Word = res.data.data;
        createUserSynonymsByWord({ input: opts.synonyms }, word);
        setWords((prev) => [...prev, res.data.data]);
      });
    },
    [makeRequest, user.id, FETCH_PARAMS, words, createUserSynonymsByWord]
  );

  return {
    updateUserWords,
    getUserWords,
    setWords,
    createUserWords,
    words,
    isLoading,
    error,
  };
}
