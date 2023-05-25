import { useState, useEffect, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Paragraph } from "../types/general";
import { UserContext } from "../context/UserContext";
import createParagraphsFactory from "../utils/createParagraphsFactory";
import updateUserParagraphsByTemplateFactory from "../utils/updateUserParagraphsByTemplateFactory";
import deleteUserParagraphsByTemplateFactory from "../utils/deleteUserParagraphsByTemplateFactory";

export default function useParagraphsController() {
  const template = useMemo((): string | undefined => {
    return window.location.pathname.split("/template/")[1];
  }, []);
  const { user } = useContext(UserContext);
  const [paragraphs, setParagraphs] = useState<Paragraph[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/template/${template}/paragraph`,
      method: "POST",
    };
  }, [user.id, template]);

  const updateUserParagraphByTemplate = useCallback(
    (paragraph: Paragraph, template: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/template/${template}/paragraph/${paragraph.id}`,
          method: "PUT",
          data: paragraph,
        },
        (res) => {
          setParagraphs(res.data.data);
        }
      );
    },
    [makeRequest, FETCH_PARAMS, user.id]
  );

  const deleteUserParagraphByWord = useCallback(
    (id: number, template: string) => {
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/template/${template}/paragraph/${id}`,
          method: "DELETE",
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS]
  );

  const deleteUserParagraphsByTemplate = useCallback(
    (values: { input: string }, template: string) => {
      const body = deleteUserParagraphsByTemplateFactory(values, paragraphs);
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/template/${template}/paragraph`,
          method: "DELETE",
          data: body,
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [FETCH_PARAMS, makeRequest, user.id, paragraphs]
  );

  const updateUserParagraphsByTemplate = useCallback(
    (values: { input: string }, template: string) => {
      const body = updateUserParagraphsByTemplateFactory(values, paragraphs);
      makeRequest(
        {
          ...FETCH_PARAMS,
          url: USER_ROUTE + `/${user.id}/template/${template}/paragraph`,
          method: "PATCH",
          data: body,
        },
        (res) => {
          setParagraphs(res.data.data);
          deleteUserParagraphsByTemplate({ ...values }, template);
        }
      );
    },
    [FETCH_PARAMS, makeRequest, user.id, paragraphs, deleteUserParagraphsByTemplate]
  );

  const getUserParagraphsByTemplate = useCallback(() => {
    makeRequest(
      {
        ...FETCH_PARAMS,
        method: "GET",
      },
      (res) => setParagraphs(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const createUserParagraphsByTemplate = useCallback(
    (opts: { paragraphs: string; template: number }) => {
      if (!opts.template) {
        return;
      }

      const paragraphs = createParagraphsFactory({
        paragraphs: opts.paragraphs,
        template_id: opts.template,
      });

      makeRequest({ ...FETCH_PARAMS, data: paragraphs }, (res) =>
        setParagraphs(res.data.data)
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
            `/${user.id}/paragraph/${[id]}?paragraphs=${[
              id,
            ]}&template=${template}`,
          method: "DELETE",
        },
        (res) => setParagraphs(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS, template]
  );

  useEffect(() => {
    if (template) {
      getUserParagraphsByTemplate();
    }
  }, [getUserParagraphsByTemplate, template]);

  return {
    updateUserParagraphByTemplate,
    setParagraphs,
    createUserParagraphsByTemplate,
    deleteParagraph,
    getUserParagraphsByTemplate,
    deleteUserParagraphByWord,
    paragraphs,
    isLoading,
    updateUserParagraphsByTemplate,
    error,
  };
}
