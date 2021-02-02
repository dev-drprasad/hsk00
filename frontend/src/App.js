import React, { useCallback, useEffect, useState } from "react";
import "./App.scss";
import Alert from "./Alert";
import GitHubIcon from "./GitHubIcon";
import { ReactComponent as Spinner } from "./spinner.svg";
import { useTranslation } from "react-i18next";
import GameList from "./GameList";
const initialState = {
  rootDir: "",
  modified: false,
  games: [],
  categoryID: -1,
  errors: { rootDir: "" },
  message: undefined,
  saving: false,
  language: "en",
};

const gameSorter = (a, b) => a.name.localeCompare(b.name);

function App() {
  const { t, i18n } = useTranslation("translation");
  const [state, setState] = useState(initialState);
  const { rootDir, games, categoryID, errors, message, modified, saving } = state;

  const handleSelectRootClick = () => {
    window.backend.Runtime.SelectRootDir().then((selectedDir) => {
      setState((s) => ({ ...s, rootDir: selectedDir }));
    });
  };

  const handleSelectGamesClick = (e) => {
    window.backend.Runtime.SelectGames()
      .then((selected) => {
        const temp = [...games, ...selected];
        console.log("temp :>> ", temp);
        let seen = {};
        const newGamesSet = [];
        temp.forEach((g) => {
          if (!(g.srcPath && seen[g.srcPath])) newGamesSet.push((seen[g.srcPath] = g));
        });
        seen = undefined;
        setState((s) => ({ ...s, modified: true, games: newGamesSet.sort(gameSorter) }));
      })
      .catch(setError);
  };

  const handleRootDirChange = (e) => {
    e.persist();
    setState((s) => ({ ...s, rootDir: e.target.value }));
  };

  const handleCategoryChange = (e) => {
    e.persist();
    setState((s) => ({ ...s, categoryID: parseInt(e.target.value) }));
  };

  const setError = (msg) => {
    setState((s) => ({
      ...s,
      saving: false,
      message: {
        type: "danger",
        content: msg,
      },
    }));
  };

  const refreshGameList = useCallback(() => {
    window.backend.Runtime.GetGameList(rootDir, categoryID)
      .then((games) => {
        console.log("gameList :>> ", games);
        setState((s) => ({
          ...s,
          modified: false,
          newGames: [],
          games: games.sort(gameSorter),
          message: undefined,
        }));
      })
      .catch(setError);
  }, [rootDir, categoryID]);

  const handleGHClick = () => {
    window.backend.Runtime.OpenURL("https://github.com/dev-drprasad/hsk00/");
  };

  const handleGameToggleDelete = (id) => () => {
    const g = { ...games[id], deleted: !games[id].deleted };
    if (g.hsk) {
      setState((s) => ({ ...s, modified: true, games: [...games.slice(0, id), g, ...games.slice(id + 1)] }));
    } else {
      // this is unsaved game, delete this from list
      setState((s) => ({ ...s, modified: true, games: [...games.slice(0, id), ...games.slice(id + 1)] }));
    }
  };

  const handleSubmit = () => {
    const errors = {};
    if (!rootDir) {
      errors["rootDir"] = "root path is empty";
    }
    if (categoryID < 0) {
      errors["categoryID"] = "select category";
    }

    if (errors.length > 0) {
      setState((s) => ({ ...s, errors: errors }));
      return;
    } else {
      setState((s) => ({ ...s, saving: true, message: undefined, errors: {} }));
    }
    console.log("rootDir, categoryID, newGames :>> ", rootDir, categoryID, games);
    window.backend.Runtime.Save(rootDir, categoryID, games)
      .then((res) => {
        setState((s) => ({
          ...s,
          games: res,
          newGames: [],
          modified: false,
          saving: false,
          message: {
            type: "success",
            content: "🎉  Changes saved successfully!  🎉",
          },
        }));
      })
      .catch(setError);
  };

  // const allgames = [...games, ...newGames].sort(gameSorter);
  const newGamesCount = games.filter(({ hsk }) => !hsk).length;
  const deletedGamesCount = games.filter(({ deleted }) => !!deleted).length;

  let totalGamesLabel = "";
  if (games.length > 0) {
    totalGamesLabel = `${games.length - newGamesCount}`;
    if (newGamesCount) totalGamesLabel += ` + ${newGamesCount}`;
    if (deletedGamesCount) totalGamesLabel += ` - ${deletedGamesCount}`;
    totalGamesLabel = `(${totalGamesLabel})`;
  }

  useEffect(() => {
    if (rootDir && categoryID > -1) {
      refreshGameList();
    }
  }, [rootDir, categoryID, refreshGameList]);

  const categoryOptions = [
    "Action Game",
    "Shoot Game",
    "Sport Game",
    "Fight Game",
    "Racing Game",
    "Puzzle Game",
  ].map((c, i) => ({ label: `${t("Category")} ${i + 1} (${c})`, value: i }));

  return (
    <React.Fragment>
      <div className="App">
        <div className="FormItem">
          <div className="label" htmlFor="rootDir">
            {t("Choose root directory")} :
          </div>
          <div className="group RootDirGroup">
            <input
              className="FormControl Input"
              name="rootDir"
              placeholder={`${t("Choose root directory")} (SD Card)`}
              onChange={handleRootDirChange}
              value={rootDir}
            />
            <button className="FormControl btn" onClick={handleSelectRootClick}>
              {t("Choose")}
            </button>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <div className="FormItem">
          <div className="label" htmlFor="rootDir">
            {t("Select game category")} :
          </div>
          <div>
            <select className="FormControl Select CategorySelect" name="categoryID" onChange={handleCategoryChange}>
              <option value={-1}>----------</option>
              {categoryOptions.map(({ label, value }) => (
                <option key={label} value={value}>
                  {label}
                </option>
              ))}
            </select>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <div className="FormItem games-list">
          <label className="label" htmlFor="rootDir">
            {t("Games")} {totalGamesLabel}:
          </label>
          <div className="list-actions">
            <button className="FormControl btn btn-sm" onClick={refreshGameList} disabled={!rootDir || categoryID === -1}>
              {t("Reset")}
            </button>
            <button
              className="FormControl btn btn-sm btn-primary"
              onClick={handleSelectGamesClick}
              disabled={!rootDir || categoryID === -1}
            >
              + {t("Add")}
            </button>
          </div>
          <GameList games={games} onToggleDelete={handleGameToggleDelete} />
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <Alert type={message?.type} message={message?.content} />
        <div className="FormItem SubmitButtonWrapper">
          <button className="FormControl SubmitButton btn btn-primary" disabled={!modified || saving} onClick={handleSubmit}>
            {t("Save Changes")}
            {saving && <Spinner className="spinner" />}
          </button>
        </div>
      </div>
      <span onClick={handleGHClick} className="github-link">
        <GitHubIcon />
      </span>
      <select className="language-select FormControl form-control-sm" onChange={(e) => i18n.changeLanguage(e.target.value)}>
        <option value="en">English</option>
        <option value="ru">русский</option>
      </select>
    </React.Fragment>
  );
}

export default App;
