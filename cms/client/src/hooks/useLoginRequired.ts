import axios from "axios";
import { useCallback, useContext, useEffect } from "react";
import { useNavigate } from "react-router";
import { USER_ROUTE } from "../constants";
import { UserContext } from "../context/UserContext";

export default function useLoginRequired() {
  const navigate = useNavigate();
  const { SetUser, user } = useContext(UserContext);

  const fetchData = useCallback(() => {
    axios
      .get(USER_ROUTE, {
        withCredentials: true,
      })
      .then((res) => {
        SetUser(res.data.data);
      })
      .catch((_) => {
        navigate("/login");
      });
  }, [navigate, SetUser]);

  useEffect(() => {
    if (!user.token || !user.email) {
      fetchData();
    }
  }, [fetchData, user.token, user.email]);
}
