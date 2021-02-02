import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import Emoji from "./Emoji";

function HasUpdate({ className }) {
  const { t } = useTranslation("translation");
  const [hasUpdate, setHasUpdate] = useState(false);

  const handleUpdateClick = () => {
    window.backend.Runtime.OpenURL("https://github.com/dev-drprasad/hsk00/releases/latest");
  };

  useEffect(() => {
    window.backend.Runtime.HasUpdate().then(setHasUpdate).catch(console.error);
  }, []);

  return (
    <a onClick={handleUpdateClick} className={className} href="https://github.com/dev-drprasad/hsk00/releases/latest">
      {hasUpdate && (
        <>
          <Emoji ariaLabel="update available">🎉</Emoji> {t("New update available!")}
        </>
      )}
    </a>
  );
}

export default HasUpdate;
