# $NetBSD: Makefile,v 1.32 2018/01/01 21:18:25 adam Exp $

DISTNAME=	holtz-1.4.0
PKGREVISION=	8
CATEGORIES=	games
MASTER_SITES=	${MASTER_SITE_SOURCEFORGE:=holtz/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://holtz.sourceforge.net/
COMMENT=	Engine for the abstract strategy games Zèrtz and Dvonn
LICENSE=	gnu-gpl-v2

GNU_CONFIGURE=		yes
USE_PKGLOCALEDIR=	yes
USE_LANGUAGES=		c c++
USE_TOOLS+=		gmake

.include "../../devel/boost-libs/buildlink3.mk"
.include "../../x11/wxGTK28/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
