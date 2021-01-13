import React from "react";
import "./Alert.scss";

function Alert({ type, message }) {
  return <div className={"alert alert-" + type}>{message}</div>;
}

export default Alert;
