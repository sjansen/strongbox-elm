module Page exposing (Page(..), view, viewErrors)

{-| Determines which navbar link (if any) will be rendered as active.
Note that we don't enumerate every page here, because the navbar doesn't
have links for every page. Anything that's not part of the navbar falls
under Other.
-}

import Browser exposing (Document)
import Html.Styled exposing (Html, a, button, div, footer, h1, header, li, main_, nav, p, text, ul)
import Html.Styled.Attributes exposing (class, classList, href, id, style)
import Html.Styled.Events exposing (onClick)
import Route exposing (Route)
import Viewer exposing (Viewer)


type Page
    = Main


{-| Take a page's Html and frames it with a header and footer.
The caller provides the current user, so we can display in either
"signed in" (rendering username) or "signed out" mode.
isLoading is for determining whether we should show a loading spinner
in the header. (This comes up during slow page transitions.)
-}
view : Maybe Viewer -> Page -> { title : String, content : Html msg } -> Document msg
view maybeViewer page { title, content } =
    { title = "Strongbox " ++ title
    , body =
        List.map
            Html.Styled.toUnstyled
            [ div [ class "grid" ]
                [ viewHeader page maybeViewer
                , viewContent content
                , viewFooter
                ]
            ]
    }


viewErrors : msg -> List String -> Html msg
viewErrors dismissErrors errors =
    if List.isEmpty errors then
        text ""

    else
        div
            [ class "error-messages"
            , style "position" "fixed"
            , style "top" "0"
            , style "background" "rgb(250, 250, 250)"
            , style "padding" "20px"
            , style "border" "1px solid"
            ]
        <|
            List.map (\error -> p [] [ text error ]) errors
                ++ [ button [ onClick dismissErrors ] [ text "Ok" ] ]



-- INTERNAL


isActive : Page -> Route -> Bool
isActive page route =
    case ( page, route ) of
        _ ->
            False


navbarLink : Page -> Route -> List (Html msg) -> Html msg
navbarLink page route linkContent =
    li [ classList [ ( "active", isActive page route ) ] ]
        [ a [ Route.href route ] linkContent ]


viewContent : Html msg -> Html msg
viewContent content =
    main_ [ id "content" ] [ content ]


viewFooter : Html msg
viewFooter =
    footer []
        [ div [] [ text "Kilroy was here." ]
        ]


viewHeader : Page -> Maybe Viewer -> Html msg
viewHeader page maybeViewer =
    let
        linkTo =
            navbarLink page
    in
    header []
        [ h1 [] [ a [ href "/" ] [ text "Strongbox" ] ]
        , nav []
            (case maybeViewer of
                Just viewer ->
                    [ ul [] [ linkTo Route.Logout [ text "Sign out" ] ] ]

                Nothing ->
                    [ ul [] [ linkTo Route.Login [ text "Sign in" ] ] ]
            )
        ]
