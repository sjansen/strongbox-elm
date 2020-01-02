module Route exposing (Route(..), fromUrl, href, replaceUrl)

import Browser.Navigation as Nav
import Html.Styled exposing (Attribute)
import Html.Styled.Attributes as Attr
import Url exposing (Url)
import Url.Parser as Parser exposing ((</>), Parser, oneOf, s)


type Route
    = Root
    | Login
    | Logout


fromUrl : Url -> Maybe Route
fromUrl url =
    url |> Parser.parse parser


href : Route -> Attribute msg
href targetRoute =
    Attr.href (routeToString targetRoute)


replaceUrl : Nav.Key -> Route -> Cmd msg
replaceUrl key route =
    Nav.replaceUrl key (routeToString route)



-- INTERNAL


parser : Parser (Route -> a) a
parser =
    oneOf
        [ Parser.map Root Parser.top
        , Parser.map Login (s "login")
        , Parser.map Logout (s "logout")
        ]


routeToPieces : Route -> List String
routeToPieces page =
    case page of
        Root ->
            []

        Login ->
            [ "login" ]

        Logout ->
            [ "logout" ]


routeToString : Route -> String
routeToString page =
    "/" ++ String.join "/" (routeToPieces page)
