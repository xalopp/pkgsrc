# $NetBSD: Makefile,v 1.17 2017/05/21 06:49:09 jlam Exp $
#

DISTNAME=		rinetd
PKGNAME=		${DISTNAME}-0.62
PKGREVISION=		3
CATEGORIES=		net
MASTER_SITES=		http://www.boutell.com/rinetd/http/

MAINTAINER=		pkgsrc-users@NetBSD.org
HOMEPAGE=		http://www.boutell.com/rinetd/
COMMENT=		Internet redirection server

NO_CONFIGURE=		yes
BUILD_TARGET=		rinetd

CPPFLAGS+=		-DRINETD_CONF="\"${PKG_SYSCONFDIR}/rinetd.conf\""

INSTALLATION_DIRS=	sbin ${PKGMANDIR}/man8

RCD_SCRIPTS=		rinetd

SUBST_CLASSES+=		paths
SUBST_FILES.paths=	rinetd.8
SUBST_SED.paths+=	-e 's,@PREFIX@,${PREFIX},g'
SUBST_SED.paths+=	-e 's,@PKG_SYSCONFDIR@,${PKG_SYSCONFDIR},g'
SUBST_STAGE.paths=	post-patch

LDFLAGS.SunOS+=		-lsocket -lnsl

do-build:
	cd ${WRKSRC} && env ${MAKE_ENV} ${CC} ${CPPFLAGS} ${CFLAGS} ${LDFLAGS} -o rinetd *.c

do-install:
	${INSTALL_PROGRAM} ${WRKSRC}/rinetd ${DESTDIR}${PREFIX}/sbin/.
	${INSTALL_MAN} ${WRKSRC}/rinetd.8 ${DESTDIR}${PREFIX}/${PKGMANDIR}/man8/.

.include "../../mk/bsd.pkg.mk"
