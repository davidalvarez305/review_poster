import { useState, useCallback, useContext, useMemo } from "react";
import { USER_ROUTE } from "../constants";
import useFetch from "./useFetch";
import { Template } from "../types/general";
import { UserContext } from "../context/UserContext";

export default function useTemplatesController() {
  const { user } = useContext(UserContext);
  const [templates, setTemplates] = useState<Template[]>([]);
  const { isLoading, makeRequest, error } = useFetch();
  const FETCH_PARAMS = useMemo(() => {
    return {
      url: USER_ROUTE + `/${user.id}/template`,
      method: "POST",
    };
  }, [user.id]);

  const getUserTemplates = useCallback(() => {
    makeRequest({ ...FETCH_PARAMS, method: "GET" }, (res) =>
      setTemplates(res.data.data)
    );
  }, [makeRequest, FETCH_PARAMS]);

  const createUserTemplates = useCallback(
    (opts: { paragraphs: string; template: string }) => {
      if (opts.template === "") {
        return;
      }
      const data = { template: opts.template, user_id: user.id };
      makeRequest({ ...FETCH_PARAMS, data }, (res) =>
        setTemplates(res.data.data)
      );
    },
    [makeRequest, user.id, FETCH_PARAMS]
  );

  return {
    getUserTemplates,
    setTemplates,
    createUserTemplates,
    templates,
    isLoading,
    error,
  };
}
