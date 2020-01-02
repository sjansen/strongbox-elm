module Viewer exposing (Viewer)

{-| The logged-in user currently viewing this page. It stores enough data to
be able to render the menu bar (username and avatar), along with Cred so it's
impossible to have a Viewer if you aren't logged in.
-}


type alias Viewer =
    { email : String
    , familyName : String
    , givenName : String
    }
