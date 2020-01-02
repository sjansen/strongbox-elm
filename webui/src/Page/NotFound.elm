module Page.NotFound exposing (view)

import Html.Styled exposing (Html, h1, text)


view : { title : String, content : Html msg }
view =
    { title = "Page Not Found"
    , content =
        h1 [] [ text "Not Found" ]
    }
