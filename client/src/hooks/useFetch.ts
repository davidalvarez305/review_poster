import { useCallback, useContext, useMemo, useState } from "react";
import axios, { AxiosRequestConfig, AxiosResponse } from "axios";
import { UserContext } from "../context/UserContext";

export default function useFetch() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState({ message: "" });
  const cancelToken = useMemo(() => axios.CancelToken.source(), []);
  const { user } = useContext(UserContext);

  const makeRequest = useCallback(
    async (
      config: AxiosRequestConfig,
      callback: (data: AxiosResponse) => void
    ) => {
      setIsLoading(true);
      setError({ message: "" });

      await axios({
        url: config.url,
        cancelToken: cancelToken.token,
        method: config.method ? config.method : undefined,
        headers: {
          ...config.headers,
          Authorization: "Bearer " + user.auth_header_string,
        },
        responseType: config.responseType ? config.responseType : undefined,
        withCredentials: true,
        data: config.data ? config.data : null,
        validateStatus: function (status) {
          return status < 400;
        },
      })
        .then((response) => {
          setIsLoading(false);
          callback(response);
        })
        .catch((error) => {
          if (error?.response?.data) {
            setError({ message: error.response.data.data });
          } else {
            setError({ message: error.message });
          }
          setIsLoading(false);
        });
    },
    [cancelToken, user.auth_header_string]
  );

  const errorCallback = (callbackMessage: string) => {
    setError({ message: callbackMessage });
    setIsLoading(false);
  };

  return {
    isLoading,
    error,
    makeRequest,
    errorCallback,
    cancelToken,
  };
}
