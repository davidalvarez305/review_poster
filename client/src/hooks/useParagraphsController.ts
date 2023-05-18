import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Paragraph, Template } from "../types/general";
import { UserContext } from "../context/UserContext";
import createParagraphsFactory from "../utils/createParagraphsFactory";
import { createUpdateParagraphsFactory } from "../utils/createUpdateParagraphsFactory";

export default function useParagraphsController() {
  const template = useMemo((): string | undefined => {
    return window.location.pathname.split("/template/")[1]
  }, []);
  const { user } = useContext(UserContext);
  const [paragraphs, setParagraphs] = useState<Paragraph[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const [editSingleParagraph, setEditSingleParagraph] =
    useState<Paragraph | null>(null);
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/paragraph`,
      method: "POST",
    };
  }, [user.id]);

  const updateParagraphs = useCallback(
    (opts: { input: string }, template_id: number, template: string) => {
      const paragraph = opts.input.split("\n")[0];
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/paragraph?template=${template}`,
          method: "PUT",
          data: {
            ...editSingleParagraph,
            name: paragraph,
            template_id,
          },
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, editSingleParagraph]
  );

  const updateParagraph = useCallback(
    (paragraph: Paragraph) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/paragraph/${paragraph.id}?template=${template}`,
          method: "PUT",
          data: paragraph,
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, template]
  );

  const bulkUpdateParagraphs = useCallback(
    (values: { input: string }) => {
      let data = createUpdateParagraphsFactory(
        paragraphs,
        values.input.split("\n"),
        paragraphs[0].template_id
      );
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/paragraph/bulk?template=${template}`,
          method: "POST",
          data,
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, paragraphs, user.id, template]
  );

  const getParagraphs = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setParagraphs(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const getParagraphsByTemplate = useCallback(() => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
        url: USER_ROUTE + `/${user.id}/paragraph?template=${template}`,
      },
      (res) => setParagraphs(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS, user.id, template]);

  const createParagraphs = useCallback(
    (opts: { paragraphs: string; template: number }) => {

      if (!opts.template) {
        return;
      }

      const paragraphs = createParagraphsFactory({
        paragraphs: opts.paragraphs,
        template_id: opts.template,
      });

      makeRequest(
        { ...FETCH_PARAMS, data: paragraphs },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS]
  );

  const deleteParagraph = useCallback(
    (id: number) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/paragraph/${[id]}?paragraphs=${[id]}&template=${template}`,
          method: "DELETE",
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS, template]
  );

  useEffect(() => {
    if (template) {
      getParagraphsByTemplate();
    } else {
      getParagraphs();
    }
  }, [getParagraphs, getParagraphsByTemplate, template]);

  return {
    updateParagraphs,
    getParagraphs,
    setParagraphs,
    createParagraphs,
    deleteParagraph,
    getParagraphsByTemplate,
    bulkUpdateParagraphs,
    setEditSingleParagraph,
    paragraphs,
    isLoading,
    updateParagraph,
    error,
  };
}
