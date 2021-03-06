# $NetBSD: Makefile,v 1.33 2016/06/01 15:00:25 richard Exp $

DISTNAME=	openal-0.0.8
PKGREVISION=	6
CATEGORIES=	audio
MASTER_SITES=	${MASTER_SITE_GENTOO:=distfiles/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://www.openal.org/
COMMENT=	3D positional spatialized sound library
LICENSE=	gnu-gpl-v2

USE_TOOLS+=		gmake makeinfo
USE_LIBTOOL=		yes
GNU_CONFIGURE=		yes

PTHREAD_AUTO_VARS=	yes

INFO_FILES=		yes
INSTALLATION_DIRS=	${PKGINFODIR} include/AL
PKGCONFIG_OVERRIDE=	admin/pkgconfig/openal.pc.in
REPLACE_SH=		admin/pkgconfig/openal-config.in
LIBS+=			${LIBOSSAUDIO}

CPPFLAGS.SunOS+=	-D__EXTENSIONS__

SUBST_CLASSES+=		fixdev
SUBST_STAGE.fixdev=	post-configure
SUBST_FILES.fixdev=	src/backends/alc_backend_bsd.c
SUBST_SED.fixdev=	-e 's;/dev/dsp;${DEVOSSSOUND};g'

.include "../../mk/bsd.prefs.mk"

.if !empty(MACHINE_PLATFORM:MDarwin-10.*-*)
CONFIGURE_ENV+=		NASM_FORMAT="-f macho"
.endif

.if ${OPSYS} == "SunOS"
.  if ${ABI:U} == "64"
CONFIGURE_ENV+=	NASM_FORMAT="-f elf64"
.  else
CONFIGURE_ENV+=	NASM_FORMAT="-f elf32"
.  endif
CONFIGURE_ENV+=	NM="${NM:Unm} -p"
LDFLAGS+=	-Wl,-zdefs
.endif

BUILD_DEPENDS+=	nasm>=0.98:../../devel/nasm

.include "../../mk/pthread.buildlink3.mk"
.include "../../mk/oss.buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
