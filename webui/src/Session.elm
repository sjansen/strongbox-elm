module Session exposing (Session, changes, fromViewer, navKey, viewer)

import Browser.Navigation as Nav
import LocalStorage
import Viewer exposing (Viewer)


type Session
    = LoggedIn Nav.Key Viewer
    | Guest Nav.Key


changes : (Session -> msg) -> Nav.Key -> Sub msg
changes toMsg key =
    LocalStorage.viewerChanges
        (\maybeViewer -> toMsg (fromViewer key maybeViewer))


fromViewer : Nav.Key -> Maybe Viewer -> Session
fromViewer key maybeViewer =
    case maybeViewer of
        Just viewerVal ->
            LoggedIn key viewerVal

        Nothing ->
            Guest key


navKey : Session -> Nav.Key
navKey session =
    case session of
        LoggedIn key _ ->
            key

        Guest key ->
            key


viewer : Session -> Maybe Viewer
viewer session =
    case session of
        LoggedIn _ val ->
            Just val

        Guest _ ->
            Nothing
