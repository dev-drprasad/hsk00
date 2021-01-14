import React, { useCallback, useEffect, useState } from "react";
import "./App.scss";
import Alert from "./Alert";
import GitHubIcon from "./GitHubIcon";
import { ReactComponent as Spinner } from "./spinner.svg";
const initialState = {
  rootDir: "",
  modified: false,
  games: [],
  newGames: [],
  categoryID: -1,
  errors: { rootDir: "" },
  message: undefined,
  saving: false,
};

const gameSorter = (a, b) => a.name.localeCompare(b.name);

function App() {
  const [state, setState] = useState(initialState);
  const { rootDir, newGames, games, categoryID, errors, message, modified, saving } = state;

  const handleSelectRootClick = () => {
    window.backend.Runtime.SelectRootDir().then((selectedDir) => {
      setState((s) => ({ ...s, rootDir: selectedDir }));
    });
  };

  const handleSelectGamesClick = (e) => {
    window.backend.Runtime.SelectGames()
      .then((selected) => {
        console.log("selected :>> ", selected);
        const temp = [...newGames, ...selected];
        let seen = {};
        const newGamesSet = [];
        temp.forEach((g) => {
          if (!seen[g.srcPath]) newGamesSet.push((seen[g.srcPath] = g));
        });
        seen = undefined;
        setState((s) => ({ ...s, modified: true, newGames: newGamesSet }));
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
          games: games,
          message: undefined,
        }));
      })
      .catch(setError);
  }, [rootDir, categoryID]);

  const handleGHClick = () => {
    window.backend.Runtime.OpenURL("https://github.com/dev-drprasad/hsk00/");
  };

  const handleSubmit = () => {
    const errors = {};
    if (!rootDir) {
      errors["rootDir"] = "root path is empty";
    }
    if (categoryID < 0) {
      errors["categoryID"] = "select category";
    }
    if (newGames.length === 0) {
      errors["newGames"] = "select game(s) to add";
    }
    if (errors.length > 0) {
      setState((s) => ({ ...s, errors: errors }));
      return;
    } else {
      setState((s) => ({ ...s, saving: true, errors: {} }));
    }
    console.log("rootDir, categoryID, newGames :>> ", rootDir, categoryID, newGames);
    window.backend.Runtime.AddGames(rootDir, categoryID, newGames)
      .then((res) => {
        console.log("res", res);
        setState((s) => ({
          ...s,
          games: res,
          newGames: [],
          modified: false,
          saving: false,
          message: {
            type: "success",
            content: `🎉  ${newGames.length} game(s) are added!  🎉`,
          },
        }));
      })
      .catch(setError);
  };

  const allgames = [...games, ...newGames].sort(gameSorter);
  console.log("allgames :>> ", allgames);
  useEffect(() => {
    if (rootDir && categoryID > -1) {
      refreshGameList();
    }
  }, [rootDir, categoryID, refreshGameList]);

  return (
    <React.Fragment>
      <div className="App">
        <div className="FormItem">
          <div className="label" htmlFor="rootDir">
            Choose root path :
          </div>
          <div className="group RootDirGroup">
            <input
              className="FormControl Input"
              name="rootDir"
              placeholder="Choose root path (SD Card)"
              onChange={handleRootDirChange}
              value={rootDir}
            />
            <button className="FormControl btn" onClick={handleSelectRootClick}>
              Choose
            </button>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <div className="FormItem">
          <div className="label" htmlFor="rootDir">
            Select game category :
          </div>
          <div>
            <select className="FormControl Select CategorySelect" name="categoryID" onChange={handleCategoryChange}>
              <option value={-1}>----------</option>
              <option value={0}>0. Action Games</option>
              <option value={1}>1. Shoot Games</option>
              <option value={2}>2. Sport Games</option>
              <option value={3}>3. Fight Games</option>
              <option value={4}>4. Racing Games</option>
              <option value={5}>5. Puzzle Games</option>
            </select>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <div className="FormItem games-list">
          <label className="label" htmlFor="rootDir">
            Games {games.length ? `(${games.length}` + (newGames.length ? ` + ${newGames.length} unsaved)` : ")") : ""}:
          </label>
          <div className="list-actions">
            <button className="FormControl btn btn-sm" onClick={refreshGameList} disabled={!rootDir || categoryID === -1}>
              Reset
            </button>
            <button
              className="FormControl btn btn-sm btn-primary"
              onClick={handleSelectGamesClick}
              disabled={!rootDir || categoryID === -1}
            >
              + Add
            </button>
          </div>
          <div className="ListBox">
            <ul role="listbox">
              {allgames.map((g) => (
                <li key={`${g.filename}${g.srcPath}${g.id}`} className={!g.hsk ? "unsaved" : ""}>
                  {g.name}
                </li>
              ))}
            </ul>
          </div>
          <span className="FormError">{errors.rootDir}</span>
        </div>
        <Alert type={message?.type} message={message?.content} />
        <div className="FormItem SubmitButtonWrapper">
          <button className="FormControl SubmitButton btn btn-primary" disabled={!modified} onClick={handleSubmit}>
            Save Changes
            {saving && <Spinner className="spinner" />}
          </button>
        </div>
      </div>
      <span onClick={handleGHClick} className="github-link">
        <GitHubIcon />
      </span>
    </React.Fragment>
  );
}

export default App;
