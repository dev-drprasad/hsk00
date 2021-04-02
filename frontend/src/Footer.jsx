import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import GitHubIcon from "./GitHubIcon";
import "./Footer.scss";

function Footer({ onLangChange }) {
  const { t } = useTranslation("translation");
  const [version, setVersion] = useState({ current: undefined, latest: undefined, hasUpdate: false });

  const handleUpdateClick = () => {
    window.backend.Runtime.OpenURL("https://github.com/dev-drprasad/hsk00/releases/latest");
  };

  const handleGHClick = () => {
    window.backend.Runtime.OpenURL("https://github.com/dev-drprasad/hsk00/");
  };

  useEffect(() => {
    window.backend.Runtime.GetVersion().then(setVersion).catch(console.error);
  }, []);

  return (
    <>
      <select className="language-select FormControl form-control-sm" onChange={(e) => onLangChange(e.target.value)}>
        <option value="en">English</option>
        <option value="ru">русский</option>
      </select>

      <span className="version">
        {version.current}{" "}
        {version.hasUpdate && (
          <span className="update-link" onClick={handleUpdateClick}>
            ({t("update available")} : {version.latest})
          </span>
        )}
      </span>
      <span onClick={handleGHClick} className="github-link">
        <GitHubIcon />
      </span>
    </>
  );
}

export default Footer;
