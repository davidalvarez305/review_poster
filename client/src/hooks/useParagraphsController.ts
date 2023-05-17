import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Paragraph, Template } from "../types/general";
import { UserContext } from "../context/UserContext";
import createParagraphsFactory from "../utils/createParagraphsFactory";
import { getId } from "../utils/getId";
import { useLocation } from "react-router-dom";
import { createUpdateParagraphsFactory } from "../utils/createUpdateParagraphsFactory";

export default function useParagraphsController() {
  const location = useLocation();
  const template = useMemo(() => location.pathname.split("/template/")[1], [location.pathname]);
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
            template_id
          },
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id, editSingleParagraph]
  );

  const bulkUpdateParagraphs = useCallback(
    (values: { input: string }) => {
      let data = createUpdateParagraphsFactory(paragraphs, values.input.split("\n"), paragraphs[0].template_id);
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

  const getParagraphsByTemplate = useCallback(
    (template: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          method: "GET",
          url: USER_ROUTE + `/${user.id}/paragraph?template=${template}`,
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, FETCH_PARAMS, user.id]
  );

  const createParagraphs = useCallback(
    (opts: { paragraphs: string; template: string }, templates: Template[]) => {
      const template_id = getId(opts.template, templates, "name");

      if (!template_id) {
        return;
      }

      const paragraphs = createParagraphsFactory({
        paragraphs: opts.paragraphs,
        template_id,
        user_id: user.id,
      });

      makeRequest(
        { ...FETCH_PARAMS, method: "POST", data: paragraphs },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS]
  );

  const deleteParagraph = useCallback(
    (id: number) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url:
            USER_ROUTE +
            `/${user.id}/paragraph?paragraphs=${[id]}&template=${template}`,
          method: "DELETE",
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS, template]
  );

  useEffect(() => {
    getParagraphs();
  }, [getParagraphs]);

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
    error,
  };
}
