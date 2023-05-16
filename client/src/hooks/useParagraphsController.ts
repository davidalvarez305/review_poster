import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Paragraph, Template, WordFormInput } from "../types/general";
import { UserContext } from "../context/UserContext";
import createParagraphsFactory from "../utils/createParagraphsFactory";
import { getId } from "../utils/getId";

export default function useParagraphsController() {
  const { user } = useContext(UserContext);
  const [paragraphs, setParagraphs] = useState<Paragraph[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/paragraph`,
      method: "POST",
    };
  }, [user.id]);

  useEffect(() => {
    makeRequest({ url: USER_ROUTE + `/${user.id}/paragraph` }, (res) =>
      setParagraphs(res.data.data)
    );
  }, [makeRequest, setParagraphs, user.id]);

  const updateParagraphs = useCallback(
    (opts: WordFormInput) => {
      const paragraphsToCreate = createParagraphsFactory(
        opts,
        user.id,
        paragraphs
      );
      makeRequest(
        { ...FETCH_PARAMS, method: "PUT", data: paragraphsToCreate },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, paragraphs]
  );

  const getParagraphs = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setParagraphs(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const createParagraphs = useCallback(
    (opts: { paragraphs: string; template: string }, templates: Template[]) => {
      const template_id = getId(opts.template, templates, "name");

      if (!template_id) {
        return;
      }

      const paragraphs = createParagraphsFactory({ paragraphs: opts.paragraphs, template_id, user_id: user.id });

      makeRequest(
        { ...FETCH_PARAMS, method: "POST", data: paragraphs },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS, paragraphs]
  );

  useEffect(() => {
    getParagraphs();
  }, [getParagraphs]);

  return {
    updateParagraphs,
    getParagraphs,
    setParagraphs,
    createParagraphs,
    paragraphs,
    isLoading,
    error,
  };
}
