import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Synonym } from "../types/general";
import { UserContext } from "../context/UserContext";
import { useLocation } from "react-router-dom";
import { createUpdateSynonymsFactory } from "../utils/createUpdateSynonymsFactory";

export default function useSynonymsController() {
  const { user } = useContext(UserContext);
  const [synonyms, setSynonyms] = useState<Synonym[]>([]);
  const [editSingleSynonym, setEditSingleSynonym] = useState<Synonym | null>(
    null
  );
  const { isLoading, makeRequest, error } = useFetch();
  const location = useLocation();
  const word = useMemo(
    () => location.pathname.split("/word/")[1],
    [location.pathname]
  );
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/synonym?word=${word}`,
      method: "POST",
    };
  }, [user.id, word]);

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
          setEditSingleSynonym(null);
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

  const deleteSynonym = useCallback(
    (id: number) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/synonym?word=${word}&synonyms=${[id]}`,
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
        return { synonym, word_id: synonyms[0].id };
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
    getSynonyms();
  }, [getSynonyms]);

  return {
    updateSynonyms,
    getSynonyms,
    setSynonyms,
    synonyms,
    isLoading,
    error,
    setEditSingleSynonym,
    editSingleSynonym,
    deleteSynonym,
    bulkUpdateSynonyms,
    word,
  };
}
